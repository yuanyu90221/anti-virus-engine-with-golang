/*
    勒索軟體（Ransomware）特徵規則集
    涵蓋：WannaCry、Ryuk、LockBit 等常見勒索軟體家族
*/

rule WannaCry_Ransomware {
    meta:
        description  = "偵測 WannaCry 勒索軟體特徵字串與加密模組標記"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://www.cisa.gov/news-events/alerts/2017/05/12/indicators-associated-wannacry-ransomware"

    strings:
        $s1 = "WannaCry" ascii wide nocase
        $s2 = "WANACRY!" ascii
        $s3 = "tasksche.exe" ascii wide
        $s4 = "msg/m_chinese(simplified).wnry" ascii wide
        $ext = ".WNCRY" ascii wide

    condition:
        2 of ($s1, $s2, $s3, $s4) or $ext
}

rule Ryuk_Ransomware {
    meta:
        description  = "偵測 Ryuk 勒索軟體特徵字串與勒索通知文字"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://www.cisa.gov/news-events/alerts/2020/10/28/ransomware-activity-targeting-healthcare-and-public-health-sector"

    strings:
        $note1 = "RyukReadMe.txt" ascii wide nocase
        $note2 = "No system is safe" ascii wide
        $note3 = "Shadow Copy" ascii wide nocase
        $marker = "RYUK" ascii wide nocase
        $cmd1 = "vssadmin Delete Shadows" ascii wide nocase
        $cmd2 = "bcdedit /set {default} recoveryenabled No" ascii wide nocase

    condition:
        $marker or (2 of ($note1, $note2, $note3)) or (any of ($cmd1, $cmd2))
}

rule LockBit_Ransomware {
    meta:
        description  = "偵測 LockBit 勒索軟體特徵字串與加密後綴"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://www.cisa.gov/news-events/cybersecurity-advisories/aa23-165a"

    strings:
        $marker1 = "LockBit" ascii wide nocase
        $ext1    = ".lockbit" ascii wide nocase
        $ext2    = ".abcd" ascii wide
        $note    = "Restore-My-Files.txt" ascii wide nocase
        $cmd1    = "vssadmin.exe delete shadows /all /quiet" ascii wide nocase

    condition:
        $marker1 or any of ($ext1, $ext2) or $note or $cmd1
}
