import * as cp from "child_process";
import * as fs from "fs";
import * as path from "path";
import * as vscode from "vscode";

type MaeveStatus = {
  SessionID: string;
  ObjectCount: number;
  TokenCount: number;
  MinScore: number;
  MaxScore: number;
  AvgScore: number;
};

type Snapshot = {
  id: string;
  name: string;
  compressedTokens: number;
  originalTokens: number;
};

type RunOptions = {
  input?: string;
  allowFailure?: boolean;
  env?: NodeJS.ProcessEnv;
};

const output = vscode.window.createOutputChannel("maeve");
const cliModule = "github.com/mugiwaraluffy56/maeve/cmd/maeve@latest";

let client: MaeveClient;
let statusProvider: StatusTreeProvider;
let snapshotProvider: SnapshotTreeProvider;
let statusBar: vscode.StatusBarItem;
let refreshTimer: NodeJS.Timeout | undefined;

export function activate(context: vscode.ExtensionContext): void {
  client = new MaeveClient(context);
  statusProvider = new StatusTreeProvider();
  snapshotProvider = new SnapshotTreeProvider();

  statusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
  statusBar.command = "maeve.refreshStatus";
  statusBar.text = "$(database) maeve";
  statusBar.tooltip = "Refresh maeve status";
  statusBar.show();

  context.subscriptions.push(
    output,
    statusBar,
    vscode.window.registerTreeDataProvider("maeve.contextLens", statusProvider),
    vscode.window.registerTreeDataProvider("maeve.snapshots", snapshotProvider),
    vscode.commands.registerCommand("maeve.initSession", initSession),
    vscode.commands.registerCommand("maeve.refreshStatus", refreshAll),
    vscode.commands.registerCommand("maeve.ingestActiveFile", ingestActiveFile),
    vscode.commands.registerCommand("maeve.ingestSelectedFiles", ingestSelectedFiles),
    vscode.commands.registerCommand("maeve.ingestGitDiff", ingestGitDiff),
    vscode.commands.registerCommand("maeve.compressContext", compressContext),
    vscode.commands.registerCommand("maeve.saveSnapshot", saveSnapshot),
    vscode.commands.registerCommand("maeve.listSnapshots", listSnapshots),
    vscode.commands.registerCommand("maeve.installCli", installCli),
    vscode.commands.registerCommand("maeve.buildCli", buildCli),
    vscode.workspace.onDidSaveTextDocument(() => {
      if (getConfig().get<boolean>("autoRefresh", true)) {
        scheduleRefresh();
      }
    }),
    vscode.workspace.onDidChangeConfiguration((event) => {
      if (event.affectsConfiguration("maeve")) {
        scheduleRefresh();
      }
    })
  );

  void refreshAll();
}

export function deactivate(): void {
  if (refreshTimer) {
    clearTimeout(refreshTimer);
  }
}

async function initSession(): Promise<void> {
  const name = await vscode.window.showInputBox({
    title: "Initialize maeve session",
    prompt: "Session name",
    value: "default",
  });
  if (name === undefined) {
    return;
  }

  await withProgress("Initializing maeve session", async () => {
    const branch = await currentGitBranch();
    const args = ["init", "--name", name.trim() || "default"];
    if (branch) {
      args.push("--branch", branch);
    }
    const result = await client.run(args);
    vscode.window.showInformationMessage(result.stdout.trim() || "maeve session initialized");
  });
  await refreshAll();
}

async function ingestActiveFile(): Promise<void> {
  const editor = vscode.window.activeTextEditor;
  if (!editor || editor.document.uri.scheme !== "file") {
    vscode.window.showWarningMessage("Open a file before ingesting it with maeve.");
    return;
  }
  await ingestFiles([editor.document.uri]);
}

async function ingestSelectedFiles(first?: vscode.Uri, rest?: vscode.Uri[]): Promise<void> {
  const files = [first, ...(rest ?? [])].filter((uri): uri is vscode.Uri => !!uri && uri.scheme === "file");
  if (files.length === 0) {
    await ingestActiveFile();
    return;
  }
  await ingestFiles(files);
}

async function ingestFiles(files: vscode.Uri[]): Promise<void> {
  await withProgress(`Ingesting ${files.length} file${files.length === 1 ? "" : "s"}`, async () => {
    for (const file of files) {
      const relative = client.relativeToWorkspace(file.fsPath);
      const result = await client.run(["ingest", "file", relative]);
      output.appendLine(result.stdout.trim());
    }
  });
  vscode.window.showInformationMessage(`maeve ingested ${files.length} file${files.length === 1 ? "" : "s"}.`);
  await refreshAll();
}

async function ingestGitDiff(): Promise<void> {
  const diff = await client.runGit(["diff", "--"]);
  if (!diff.stdout.trim()) {
    vscode.window.showInformationMessage("No Git diff to ingest.");
    return;
  }

  await withProgress("Ingesting Git diff", async () => {
    const result = await client.run(["ingest", "diff"], { input: diff.stdout });
    output.appendLine(result.stdout.trim());
  });
  vscode.window.showInformationMessage("maeve ingested the current Git diff.");
  await refreshAll();
}

async function compressContext(): Promise<void> {
  const budget = getConfig().get<number>("tokenBudget", 8000);
  const result = await withProgress("Compressing maeve context", () => client.run(["compress", "--budget", String(budget)]));
  const payload = result.stdout.trim();
  if (!payload) {
    vscode.window.showInformationMessage("maeve produced an empty compressed context.");
    return;
  }

  const doc = await vscode.workspace.openTextDocument({
    language: "markdown",
    content: payload,
  });
  await vscode.window.showTextDocument(doc, { preview: false });
  await vscode.env.clipboard.writeText(payload);
  vscode.window.showInformationMessage("Compressed maeve context opened and copied to clipboard.");
}

async function saveSnapshot(): Promise<void> {
  const name = await vscode.window.showInputBox({
    title: "Save maeve snapshot",
    prompt: "Snapshot name",
    value: `snapshot ${new Date().toISOString().slice(0, 19).replace("T", " ")}`,
  });
  if (name === undefined) {
    return;
  }

  const budget = getConfig().get<number>("tokenBudget", 8000);
  const result = await withProgress("Saving maeve snapshot", () =>
    client.run(["snapshot", "save", name.trim() || "snapshot", "--budget", String(budget)])
  );
  vscode.window.showInformationMessage(result.stdout.trim() || "maeve snapshot saved");
  await refreshAll();
}

async function listSnapshots(): Promise<void> {
  await refreshSnapshots();
  await vscode.commands.executeCommand("maeve.snapshots.focus");
}

async function buildCli(): Promise<void> {
  const workspace = client.workspaceFolder();
  if (!workspace) {
    vscode.window.showWarningMessage("Open a workspace before building the maeve CLI.");
    return;
  }

  const binDir = path.join(workspace.uri.fsPath, ".maeve", "bin");
  const binary = path.join(binDir, process.platform === "win32" ? "maeve.exe" : "maeve");
  await fs.promises.mkdir(binDir, { recursive: true });

  await withProgress("Building maeve CLI", () => client.runTool("go", ["build", "-o", binary, "./cmd/maeve"]));
  vscode.window.showInformationMessage(`Built maeve CLI at ${binary}`);
  await refreshAll();
}

async function installCli(): Promise<void> {
  const binDir = client.extensionBinDir();
  const binary = client.extensionBinaryPath();
  await fs.promises.mkdir(binDir, { recursive: true });

  await withProgress("Installing maeve CLI", () =>
    client.runTool("go", ["install", cliModule], {
      env: {
        ...process.env,
        GOBIN: binDir,
      },
    })
  );
  vscode.window.showInformationMessage(`Installed maeve CLI at ${binary}`);
  await refreshAll();
}

async function refreshAll(): Promise<void> {
  await Promise.all([refreshStatus(), refreshSnapshots()]);
}

async function refreshStatus(): Promise<void> {
  try {
    const status = await client.status();
    statusProvider.setStatus(status);
    updateStatusBar(status);
  } catch (error) {
    statusProvider.setError(errorMessage(error));
    statusBar.text = "$(warning) maeve";
    statusBar.tooltip = `${errorMessage(error)}\nRun maeve: Install CLI, maeve: Build Workspace CLI, or set maeve.executablePath.`;
  }
}

async function refreshSnapshots(): Promise<void> {
  try {
    const snapshots = await client.snapshots();
    snapshotProvider.setSnapshots(snapshots);
  } catch (error) {
    snapshotProvider.setError(errorMessage(error));
  }
}

function scheduleRefresh(): void {
  if (refreshTimer) {
    clearTimeout(refreshTimer);
  }
  refreshTimer = setTimeout(() => void refreshAll(), 500);
}

function updateStatusBar(status: MaeveStatus): void {
  const budget = getConfig().get<number>("tokenBudget", 8000);
  statusBar.text = `$(database) maeve ${formatNumber(status.TokenCount)} / ${formatNumber(budget)}`;
  statusBar.tooltip = [
    `Session: ${status.SessionID}`,
    `Objects: ${status.ObjectCount}`,
    `Tokens: ${status.TokenCount}`,
    `Average score: ${status.AvgScore.toFixed(2)}`,
    `Max score: ${status.MaxScore.toFixed(2)}`,
  ].join("\n");
}

async function currentGitBranch(): Promise<string> {
  try {
    const result = await client.runGit(["branch", "--show-current"], { allowFailure: true });
    return result.stdout.trim();
  } catch {
    return "";
  }
}

async function withProgress<T>(title: string, task: () => Promise<T>): Promise<T> {
  return vscode.window.withProgress(
    {
      location: vscode.ProgressLocation.Notification,
      title,
      cancellable: false,
    },
    task
  );
}

function getConfig(): vscode.WorkspaceConfiguration {
  return vscode.workspace.getConfiguration("maeve");
}

function formatNumber(value: number): string {
  return new Intl.NumberFormat("en", { maximumFractionDigits: 0 }).format(value);
}

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : String(error);
}

class MaeveClient {
  constructor(private readonly context: vscode.ExtensionContext) {}

  workspaceFolder(): vscode.WorkspaceFolder | undefined {
    const folders = vscode.workspace.workspaceFolders;
    return folders && folders.length > 0 ? folders[0] : undefined;
  }

  relativeToWorkspace(filePath: string): string {
    const workspace = this.workspaceFolder();
    if (!workspace) {
      return filePath;
    }
    const relative = path.relative(workspace.uri.fsPath, filePath);
    return relative && !relative.startsWith("..") ? relative : filePath;
  }

  async status(): Promise<MaeveStatus> {
    const result = await this.run(["status", "--output", "json"]);
    return JSON.parse(result.stdout) as MaeveStatus;
  }

  async snapshots(): Promise<Snapshot[]> {
    const result = await this.run(["snapshot", "list"]);
    return result.stdout
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter(Boolean)
      .map(parseSnapshotLine);
  }

  run(args: string[], options: RunOptions = {}): Promise<{ stdout: string; stderr: string }> {
    return this.runTool(this.executablePath(), [...this.globalArgs(), ...args], options);
  }

  runGit(args: string[], options: RunOptions = {}): Promise<{ stdout: string; stderr: string }> {
    return this.runTool("git", args, options);
  }

  runTool(command: string, args: string[], options: RunOptions = {}): Promise<{ stdout: string; stderr: string }> {
    const workspace = this.workspaceFolder();
    const cwd = workspace?.uri.fsPath ?? process.cwd();
    output.appendLine(`$ ${command} ${args.join(" ")}`);

    return new Promise((resolve, reject) => {
      const child = cp.spawn(command, args, {
        cwd,
        env: options.env,
        shell: false,
        windowsHide: true,
      });
      let stdout = "";
      let stderr = "";

      child.stdout.on("data", (chunk: Buffer) => {
        stdout += chunk.toString();
      });
      child.stderr.on("data", (chunk: Buffer) => {
        stderr += chunk.toString();
      });
      child.on("error", (error) => {
        reject(new Error(`${command} failed to start: ${error.message}`));
      });
      child.on("close", (code) => {
        if (stdout.trim()) {
          output.appendLine(stdout.trim());
        }
        if (stderr.trim()) {
          output.appendLine(stderr.trim());
        }
        if (code === 0 || options.allowFailure) {
          resolve({ stdout, stderr });
          return;
        }
        reject(new Error(stderr.trim() || `${command} exited with code ${code}`));
      });

      if (options.input) {
        child.stdin.write(options.input);
      }
      child.stdin.end();
    });
  }

  extensionBinDir(): string {
    return path.join(this.context.globalStorageUri.fsPath, "bin");
  }

  extensionBinaryPath(): string {
    return path.join(this.extensionBinDir(), process.platform === "win32" ? "maeve.exe" : "maeve");
  }

  private executablePath(): string {
    const workspace = this.workspaceFolder();
    if (workspace) {
      const localBinary = path.join(
        workspace.uri.fsPath,
        ".maeve",
        "bin",
        process.platform === "win32" ? "maeve.exe" : "maeve"
      );
      if (fs.existsSync(localBinary)) {
        return localBinary;
      }
    }
    const extensionBinary = this.extensionBinaryPath();
    if (fs.existsSync(extensionBinary)) {
      return extensionBinary;
    }
    return getConfig().get<string>("executablePath", "maeve");
  }

  private globalArgs(): string[] {
    const args: string[] = [];
    const configPath = getConfig().get<string>("configPath", "");
    const storePath = getConfig().get<string>("storePath", "");
    if (configPath) {
      args.push("--config", configPath);
    }
    if (storePath) {
      args.push("--store", storePath);
    }
    return args;
  }
}

function parseSnapshotLine(line: string): Snapshot {
  const parts = line.split("\t");
  const id = parts[0] ?? "";
  const name = parts[1] ?? "snapshot";
  const tokenPart = parts[2] ?? "0/0 tokens";
  const match = tokenPart.match(/(\d+)\/(\d+)/);
  return {
    id,
    name,
    compressedTokens: match ? Number(match[1]) : 0,
    originalTokens: match ? Number(match[2]) : 0,
  };
}

class StatusTreeProvider implements vscode.TreeDataProvider<vscode.TreeItem> {
  private readonly emitter = new vscode.EventEmitter<void>();
  readonly onDidChangeTreeData = this.emitter.event;
  private status: MaeveStatus | undefined;
  private error = "";

  setStatus(status: MaeveStatus): void {
    this.status = status;
    this.error = "";
    this.emitter.fire();
  }

  setError(message: string): void {
    this.status = undefined;
    this.error = message;
    this.emitter.fire();
  }

  getTreeItem(element: vscode.TreeItem): vscode.TreeItem {
    return element;
  }

  getChildren(): vscode.TreeItem[] {
    if (this.error) {
      const canBuild = fs.existsSync(path.join(client.workspaceFolder()?.uri.fsPath ?? "", "cmd", "maeve"));
      return [
        treeItem("CLI setup needed", this.error, "warning"),
        commandItem("Install CLI", "maeve.installCli", "cloud-download"),
        ...(canBuild ? [commandItem("Build workspace CLI", "maeve.buildCli", "tools")] : []),
        commandItem("Refresh", "maeve.refreshStatus", "refresh"),
      ];
    }
    if (!this.status) {
      return [treeItem("No status loaded", "Run refresh to load maeve status.", "info")];
    }

    const status = this.status;
    return [
      treeItem("Session", status.SessionID, "database"),
      treeItem("Objects", String(status.ObjectCount), "symbol-array"),
      treeItem("Tokens", String(status.TokenCount), "symbol-number"),
      treeItem("Average score", status.AvgScore.toFixed(2), "pulse"),
      treeItem("Max score", status.MaxScore.toFixed(2), "flame"),
      commandItem("Initialize session", "maeve.initSession", "database"),
      commandItem("Ingest active file", "maeve.ingestActiveFile", "file-add"),
      commandItem("Ingest Git diff", "maeve.ingestGitDiff", "diff"),
      commandItem("Compress context", "maeve.compressContext", "archive"),
      commandItem("Save snapshot", "maeve.saveSnapshot", "save"),
    ];
  }
}

class SnapshotTreeProvider implements vscode.TreeDataProvider<vscode.TreeItem> {
  private readonly emitter = new vscode.EventEmitter<void>();
  readonly onDidChangeTreeData = this.emitter.event;
  private snapshots: Snapshot[] = [];
  private error = "";

  setSnapshots(snapshots: Snapshot[]): void {
    this.snapshots = snapshots;
    this.error = "";
    this.emitter.fire();
  }

  setError(message: string): void {
    this.snapshots = [];
    this.error = message;
    this.emitter.fire();
  }

  getTreeItem(element: vscode.TreeItem): vscode.TreeItem {
    return element;
  }

  getChildren(): vscode.TreeItem[] {
    if (this.error) {
      return [treeItem("Snapshots unavailable", this.error, "warning")];
    }
    if (this.snapshots.length === 0) {
      return [
        treeItem("No snapshots", "Save a snapshot to pin compressed context.", "history"),
        commandItem("Save snapshot", "maeve.saveSnapshot", "save"),
      ];
    }
    return this.snapshots.map((snapshot) => {
      const item = new vscode.TreeItem(snapshot.name, vscode.TreeItemCollapsibleState.None);
      item.description = `${snapshot.compressedTokens}/${snapshot.originalTokens} tokens`;
      item.tooltip = snapshot.id;
      item.iconPath = new vscode.ThemeIcon("history");
      return item;
    });
  }
}

function treeItem(label: string, description: string, icon: string): vscode.TreeItem {
  const item = new vscode.TreeItem(label, vscode.TreeItemCollapsibleState.None);
  item.description = description;
  item.tooltip = description;
  item.iconPath = new vscode.ThemeIcon(icon);
  return item;
}

function commandItem(label: string, command: string, icon: string): vscode.TreeItem {
  const item = new vscode.TreeItem(label, vscode.TreeItemCollapsibleState.None);
  item.command = { command, title: label };
  item.iconPath = new vscode.ThemeIcon(icon);
  return item;
}
