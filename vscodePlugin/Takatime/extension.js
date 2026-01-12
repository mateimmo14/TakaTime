const vscode = require("vscode");
const env = require("./Plugin/StatusBarUpdate");
const fs = require("fs");
const path = require("path");
const os = require("os");

/**
 * @param {vscode.ExtensionContext} context
 */
async function activate(context) {
  console.log("TakaTime: Initializing...");

  // 1. Create Status Bar
  const statusBar = vscode.window.createStatusBarItem(
    vscode.StatusBarAlignment.Left,
    100
  );
  statusBar.text = "$(sync~spin) TakaTime: Checking...";
  statusBar.command = "takatime.setup"; // 👈 Clicking this runs the setup command
  statusBar.show();
  context.subscriptions.push(statusBar);

  // 2. Register the Setup Command
  const setupCommand = vscode.commands.registerCommand(
    "takatime.setup",
    async () => {
      // Ask user for URI
      const uri = await vscode.window.showInputBox({
        placeHolder: "mongodb+srv://admin:password@...",
        prompt: "Enter your MongoDB Connection String to start tracking",
        ignoreFocusOut: true,
        password: true, // Hides the text for privacy
      });

      if (!uri) return; // User cancelled

      // Read existing config to preserve VERSION
      const homeDir = os.homedir();
      const configPath = path.join(homeDir, ".takatime.json");
      let currentConfig = {};

      try {
        if (fs.existsSync(configPath)) {
          currentConfig = JSON.parse(fs.readFileSync(configPath, "utf8"));
        }
      } catch (e) {
        /* ignore */
      }

      // Update URI
      currentConfig.MONGO_URI = uri;
      if (!currentConfig.VERSION) currentConfig.VERSION = "v1.0.0"; // Ensure version exists

      // Save back to file
      try {
        fs.writeFileSync(configPath, JSON.stringify(currentConfig, null, 4));
        vscode.window.showInformationMessage("TakaTime: Configuration Saved! ");

        // 🔄 RE-RUN CHECKS immediately
        env.checkStatus(statusBar);
      } catch (err) {
        vscode.window.showErrorMessage(`Failed to save config: ${err.message}`);
      }
    }
  );

  context.subscriptions.push(setupCommand);

  // 3. Initial Check on Startup
  env.checkStatus(statusBar);
}

function deactivate() {}

module.exports = {
  activate,
  deactivate,
};
