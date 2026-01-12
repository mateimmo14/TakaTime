const vscode = require("vscode");
const statusHelper = require("./Plugin/StatusBarUpdate");
const env = require("./Plugin/Config");
const downloader = require("./Plugin/BinaryDownload"); // Import the new file
const fs = require("fs");
const path = require("path");
const os = require("os");

/**
 * @param {vscode.ExtensionContext} context
 */
async function activate(context) {
  console.log("TakaTime: Initializing...");

  const statusBar = vscode.window.createStatusBarItem(
    vscode.StatusBarAlignment.Left,
    100
  );
  statusBar.text = "$(sync~spin) TakaTime: Checking...";
  statusBar.command = "takatime.setup";
  statusBar.show();
  context.subscriptions.push(statusBar);

  // --- SMART SETUP COMMAND ---
  const setupCommand = vscode.commands.registerCommand(
    "takatime.setup",
    async () => {
      const config = env.getConfig();

      // CASE 1: Config/URI is missing -> Ask for it
      if (!config || !config.MONGO_URI) {
        const uri = await vscode.window.showInputBox({
          placeHolder: "mongodb+srv://admin:password@...",
          prompt: "Enter your MongoDB Connection String",
          ignoreFocusOut: true,
          password: true,
        });

        if (!uri) return;

        // Save Config Logic (Same as before)
        const homeDir = os.homedir();
        const configPath = path.join(homeDir, ".takatime.json");
        let newConfig = config || { VERSION: env.CURRENT_VERSION };
        newConfig.MONGO_URI = uri;
        if (!newConfig.VERSION) newConfig.VERSION = env.CURRENT_VERSION;

        try {
          fs.writeFileSync(configPath, JSON.stringify(newConfig, null, 4));
          vscode.window.showInformationMessage("Configuration Saved!");
          // Check status again (will likely prompt for download next click)
          statusHelper.checkStatus(statusBar);
        } catch (e) {
          vscode.window.showErrorMessage("Failed to save config");
        }
        return;
      }

      // CASE 2: Binary is Missing -> Download it
      const isBinaryReady = env.checkBinary(env.CURRENT_VERSION);
      if (!isBinaryReady) {
        try {
          const success = await downloader.downloadBinary(env.CURRENT_VERSION);
          if (success) {
            vscode.window.showInformationMessage(
              `TakaTime ${env.CURRENT_VERSION} installed successfully! `
            );
            statusHelper.checkStatus(statusBar); // Should turn Green now
          }
        } catch (err) {
          vscode.window.showErrorMessage(`Download Failed: ${err.message}`);
        }
        return;
      }

      // CASE 3: Everything is good
      vscode.window.showInformationMessage("TakaTime is active and running! ");
    }
  );

  context.subscriptions.push(setupCommand);
  statusHelper.checkStatus(statusBar);
}

function deactivate() {}

module.exports = { activate, deactivate };
