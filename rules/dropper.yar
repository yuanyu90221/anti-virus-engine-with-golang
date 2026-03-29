/*
    植入程式與下載器（Dropper / Downloader）特徵規則集
    偵測將惡意酬載寫入磁碟或從遠端下載並執行的行為特徵
*/

rule Generic_Dropper_Behavior {
    meta:
        description  = "偵測常見 dropper 行為：解碼並寫入磁碟後執行"
        severity     = "medium"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1059/"

    strings:
        $dec1 = "FromBase64String" ascii wide nocase
        $dec2 = "base64_decode" ascii wide nocase
        $write = "WriteAllBytes" ascii wide nocase
        $exec1 = "Process.Start" ascii wide nocase
        $exec2 = "ShellExecute" ascii wide nocase
        $tmp   = "%TEMP%" ascii wide nocase

    condition:
        ($dec1 or $dec2) and ($write or $exec1 or $exec2) and $tmp
}

rule PowerShell_Downloader {
    meta:
        description  = "偵測 PowerShell 下載並執行酬載的典型手法"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1105/"

    strings:
        $dl1 = "DownloadFile" ascii wide nocase
        $dl2 = "DownloadString" ascii wide nocase
        $dl3 = "WebClient" ascii wide nocase
        $iex = "IEX" ascii wide
        $iex2 = "Invoke-Expression" ascii wide nocase
        $enc = "-EncodedCommand" ascii wide nocase
        $bypass = "ExecutionPolicy Bypass" ascii wide nocase

    condition:
        ($dl1 or $dl2 or $dl3) and ($iex or $iex2 or $enc or $bypass)
}

rule Mshta_Dropper {
    meta:
        description  = "偵測透過 mshta.exe 執行遠端 HTA 酬載"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1218/005/"

    strings:
        $mshta = "mshta.exe" ascii wide nocase
        $http  = "http://" ascii wide nocase
        $https = "https://" ascii wide nocase
        $vbs   = "vbscript:" ascii wide nocase

    condition:
        $mshta and ($http or $https or $vbs)
}
