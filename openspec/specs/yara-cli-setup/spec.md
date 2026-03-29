## Requirements

### Requirement: 支援透過 apt 在 Debian/Ubuntu 安裝 YARA
在 Debian 或 Ubuntu 系統上，使用者應能透過系統套件管理工具安裝 YARA CLI。安裝完成後，`yara` binary 應位於系統 PATH 中且可直接執行。

#### Scenario: 透過 apt 安裝 YARA 並驗證版本
- **WHEN** 執行以下指令：
  ```
  sudo apt-get update
  sudo apt-get install -y yara
  ```
- **THEN** 執行 `yara --version` 應輸出版本字串（例如 `4.x.x`）且結束碼為 0

### Requirement: 支援透過 Homebrew 在 macOS 安裝 YARA
在 macOS 系統上，使用者應能透過 Homebrew 安裝 YARA CLI。安裝完成後，`yara` binary 應位於系統 PATH 中且可直接執行。

#### Scenario: 透過 Homebrew 安裝 YARA 並驗證版本
- **WHEN** 執行以下指令：
  ```
  brew install yara
  ```
- **THEN** 執行 `yara --version` 應輸出版本字串且結束碼為 0

### Requirement: 支援從原始碼編譯 YARA（適用於無套件管理工具的環境）
在無法使用套件管理工具的環境（例如自訂 CI 映像、Alpine Linux）中，使用者應能從官方 GitHub 原始碼編譯並安裝 YARA。

#### Scenario: 從原始碼編譯並安裝 YARA
- **WHEN** 依序執行以下步驟：
  ```
  # 安裝編譯依賴
  sudo apt-get install -y automake libtool make gcc pkg-config \
      libssl-dev libjansson-dev libmagic-dev

  # 取得原始碼（以 v4.5.2 為例）
  curl -L https://github.com/VirusTotal/yara/archive/refs/tags/v4.5.2.tar.gz \
      | tar xz
  cd yara-4.5.2
  ./bootstrap.sh
  ./configure --with-crypto --enable-magic --enable-dotnet
  make
  sudo make install
  ```
- **THEN** 執行 `yara --version` 應輸出 `4.5.2` 且結束碼為 0

### Requirement: 驗證 YARA 安裝可正確執行規則比對
安裝完成後，使用者應能透過一個最小化的規則與測試檔案，驗證 YARA 安裝正常且可用於防毒引擎。

#### Scenario: 以最小規則驗證 YARA 安裝
- **WHEN** 建立以下測試規則檔案 `test.yar`：
  ```yara
  rule HelloWorld {
      strings:
          $a = "Hello, YARA"
      condition:
          $a
  }
  ```
  並建立包含目標字串的測試檔案 `sample.txt`（內容含 `Hello, YARA`），
  然後執行 `yara test.yar sample.txt`
- **THEN** stdout 輸出包含 `HelloWorld sample.txt` 且結束碼為 0

#### Scenario: 無比對時 YARA 以 exit 1 結束
- **WHEN** 對不含目標字串的檔案執行 `yara test.yar other.txt`
- **THEN** 無 stdout 輸出且結束碼為 1

### Requirement: YARA 版本相容性要求
防毒引擎應與 YARA 4.x 版本相容。使用者安裝的 YARA binary 版本應 >= 4.0.0。不支援 YARA 3.x 或更舊版本。

#### Scenario: 驗證已安裝版本符合最低需求
- **WHEN** 執行 `yara --version`
- **THEN** 輸出的主版本號應 >= 4
