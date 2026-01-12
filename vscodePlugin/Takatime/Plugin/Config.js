const vscode = require("vscode");
const fs = require("fs");
const path = require("path");
const os = require("os");

const CURRENT_VERSION = "v2.0.3";

function getConfig() {
  const homeDir = os.homedir();
  const configPath = path.join(homeDir, ".takatime.json");

  // 1. Auto-Create if missing 🛠️
  if (!fs.existsSync(configPath)) {
    const defaultConfig = {
      MONGO_URI: "",
      VERSION: CURRENT_VERSION, // Default version to match your GitHub Release
    };

    try {
      fs.writeFileSync(configPath, JSON.stringify(defaultConfig, null, 4));

      // UX: Open the file for them immediately!
      vscode.workspace.openTextDocument(configPath).then((doc) => {
        vscode.window.showTextDocument(doc);
      });

      vscode.window.showInformationMessage(
        `TakaTime: Created config at ${configPath}. Please add your MONGO_URI.`
      );
      return null; // Stop here, they need to edit the file
    } catch (err) {
      vscode.window.showErrorMessage(
        `TakaTime: Failed to create config: ${err.message}`
      );
      return null;
    }
  }

  // 2. Read and Parse 📖
  try {
    const rawConfig = fs.readFileSync(configPath, "utf8");
    const config = JSON.parse(rawConfig);

    if (!config.MONGO_URI) {
      vscode.window.showWarningMessage(
        `TakaTime: MONGO_URI is empty in .takatime.json.`
      );
      return null;
    }

    // Default to a version if they deleted it
    if (!config.VERSION) {
      config.VERSION = "v1.0.0";
    }

    console.log(`✅ Config loaded. Version: ${config.VERSION}`);
    return config; // Return the full object
  } catch (err) {
    vscode.window.showErrorMessage(
      `TakaTime: Invalid JSON in config. ${err.message}`
    );
    return null;
  }
}

// We accept 'version' here now, which we will use later for downloading
function checkBinary(version) {
  const homeDir = os.homedir();
  let binName = null;
  if (process.platform === "linux") {
    binName = "taka-uploader";
  } else if (process.platform === "darwin") {
    binName = "taka-uploader";
  } else if (process.platform === "win32") {
    binName = "taka-uploader.exe";
  } else {
    vscode.window.showErrorMessage(
      `TakaTime: Unsupported platform: ${process.platform}`
    );
    return false;
  }

  // 👇 WE ADD VERSION TO THE LOCAL FILENAME HERE
  if (process.platform === "win32") {
    binName = `taka-uploader-${version}.exe`;
  } else {
    binName = `taka-uploader-${version}`;
  }

  const binaryPath = path.join(homeDir, ".takatime", "bin", binName);

  if (!fs.existsSync(binaryPath)) {
    // We can now use the version in the warning!
    vscode.window.showWarningMessage(
      `TakaTime: Binary ${version} missing. Auto-download needed.`
    );
    return false;
  }
  console.log("✅ Binary found.");
  return true;
}

module.exports = {
  getConfig,
  checkBinary,
  CURRENT_VERSION,
};
