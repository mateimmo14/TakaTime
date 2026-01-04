## 🚀 Release: v2.0.3-beta (The "User Experience" Update)

**Description:** This is a massive update focused on User Experience and Ease of Installation. We have completely removed the need to hardcode secrets in your Lua config. 
The plugin now handles binary management automatically, 
downloading the correct tools for your OS upon installation.

**Changelog:**
- feat(core): add :TakaInit command for secure, interactive setup.
- feat(install): implement auto-download logic for taka-upload and taka-report binaries.
- feat(storage): move secret storage to stdpath("data") (secure JSON file) instead of init.lua.
- fix(ci): update release workflow to build and upload taka-report binary (fixes "asset not found" error).
- fix(ui): silence "Syncing..." messages by default (set debug=false in config).
- refactor: split logic into core, storage, and utils modules for better maintainability.
  
----

## 📦 Release: v1.0.1 (The "Foundation" Release)

**Description:** The first stable release of TakaTime.nvim. This version lays the groundwork for privacy-focused, self-hosted time tracking. It connects Neovim directly to your MongoDB instance using a 
high-performance Go binary.

**Chnage logs**
- feat: initial release of Lua plugin structure.
- feat: implement Go binary for MongoDB uploads.
- feat: add debounce logic to prevent spamming the database on every keystroke.
- config: basic setup function with mongo_uri configuration support.
