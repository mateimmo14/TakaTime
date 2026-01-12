const fs = require("fs");
const path = require("path");
const os = require("os");
const https = require("https");
const vscode = require("vscode");

function getPlatformFilename() {
  const plat = process.platform; // 'win32', 'linux', 'darwin'
  const arch = process.arch; // 'x64', 'arm64'

  let osStr = "";
  if (plat === "win32") osStr = "windows";
  else if (plat === "linux") osStr = "linux";
  else if (plat === "darwin") osStr = "darwin";
  else return null;

  let archStr = "";
  if (arch === "x64") archStr = "amd64";
  else if (arch === "arm64") archStr = "arm64";
  // M1/M2 Mac support, or fallback
  else return null;

  const ext = plat === "win32" ? ".exe" : "";

  // OUTPUT EXAMPLE: taka-upload-windows-amd64.exe
  return `taka-upload-${osStr}-${archStr}${ext}`;
}

async function downloadBinary(version) {
  const filename = getPlatformFilename();
  if (!filename) {
    vscode.window.showErrorMessage(
      `TakaTime: Unsupported Platform (${process.platform}-${process.arch})`
    );
    return false;
  }

  const ext = process.platform === "win32" ? ".exe" : "";
  const localFilename = `taka-uploader-${version}${ext}`; // e.g. taka-uploader-v2.0.4.exe

  const homeDir = os.homedir();
  const binDir = path.join(homeDir, ".takatime", "bin");
  const destPath = path.join(binDir, localFilename); // 👈 Save to this path

  // Construct the URL
  // e.g. https://github.com/Rtarun3606k/TakaTime/releases/download/v2.0.4/taka-upload-windows-amd64.exe
  const url = `https://github.com/Rtarun3606k/TakaTime/releases/download/${version}/${filename}`;

  // Ensure directory exists
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Show Progress Bar in VS Code
  return vscode.window.withProgress(
    {
      location: vscode.ProgressLocation.Notification,
      title: `Downloading TakaTime ${version}...`,
      cancellable: false,
    },
    async (progress) => {
      return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(destPath);

        // Function to handle redirects (GitHub uses redirects for downloads)
        const request = (uri) => {
          https
            .get(uri, (response) => {
              // Handle Redirects (301, 302)
              if (response.statusCode === 301 || response.statusCode === 302) {
                return request(response.headers.location);
              }

              if (response.statusCode !== 200) {
                reject(
                  new Error(`Download failed: HTTP ${response.statusCode}`)
                );
                return;
              }

              response.pipe(file);

              file.on("finish", () => {
                file.close(() => {
                  // Make executable on Linux/Mac
                  if (process.platform !== "win32") {
                    try {
                      fs.chmodSync(destPath, 0o755);
                    } catch (e) {
                      console.error("Chmod failed", e);
                    }
                  }
                  resolve(true);
                });
              });
            })
            .on("error", (err) => {
              fs.unlink(destPath, () => {}); // Delete partial file
              reject(err);
            });
        };

        request(url);
      });
    }
  );
}

module.exports = { downloadBinary };
