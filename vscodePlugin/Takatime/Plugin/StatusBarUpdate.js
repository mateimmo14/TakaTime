// Helper function to update the bar based on state
const vscode = require("vscode");
const env = require("./Config");

function checkStatus(statusBar) {
  try {
    const config = env.getConfig();

    if (!config || !config.MONGO_URI) {
      statusBar.text = "$(alert) TakaTime: Setup Needed";
      statusBar.tooltip = "Click to configure MongoDB URI";
      // Highlight with a warning color
      statusBar.backgroundColor = new vscode.ThemeColor(
        "statusBarItem.warningBackground"
      );
      return;
    }

    const isBinaryReady = env.checkBinary(config.VERSION);
    if (!isBinaryReady) {
      statusBar.text = "$(cloud-download) TakaTime: Binary Missing";
      statusBar.tooltip = "Binary missing. Auto-download coming soon.";
      statusBar.backgroundColor = undefined; // Reset color
      return;
    }

    // Success State
    statusBar.text = "$(check) TakaTime: Active";
    statusBar.tooltip = `Tracking to: ${config.MONGO_URI.substring(0, 15)}...`;
    statusBar.backgroundColor = undefined;
  } catch (err) {
    console.error(err);
    statusBar.text = "$(error) TakaTime: Error";
    statusBar.tooltip = err.message;
  }
}

module.exports = {
  checkStatus,
};
