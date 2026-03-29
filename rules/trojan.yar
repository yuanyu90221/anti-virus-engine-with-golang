/*
    木馬與遠端存取工具（Trojan / RAT）特徵規則集
    涵蓋：常見後門、RAT 框架的識別字串
*/

rule Generic_RAT_Strings {
    meta:
        description  = "偵測常見 RAT 框架共用的識別字串與功能標記"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/tactics/TA0011/"

    strings:
        $s1 = "keylogger" ascii wide nocase
        $s2 = "screenshot" ascii wide nocase
        $s3 = "reverse shell" ascii wide nocase
        $s4 = "bind shell" ascii wide nocase
        $s5 = "RemoteAdmin" ascii wide nocase
        $cmd = "cmd.exe /c" ascii wide nocase

    condition:
        3 of them
}

rule NjRAT_Trojan {
    meta:
        description  = "偵測 njRAT（Bladabindi）木馬特徵字串"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://malpedia.caad.fkie.fraunhofer.de/details/win.njrat"

    strings:
        $id1 = "njRAT" ascii wide nocase
        $id2 = "Bladabindi" ascii wide nocase
        $mutex = "njq8" ascii wide
        $reg   = "SOFTWARE\\njRAT" ascii wide nocase

    condition:
        any of them
}

rule AsyncRAT_Trojan {
    meta:
        description  = "偵測 AsyncRAT 遠端存取工具特徵字串"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://malpedia.caad.fkie.fraunhofer.de/details/win.asyncrat"

    strings:
        $id1 = "AsyncClient" ascii wide nocase
        $id2 = "AsyncRAT" ascii wide nocase
        $cfg = "Pastebin" ascii wide nocase
        $key = "HKCU\\Software\\AsyncRAT" ascii wide nocase

    condition:
        any of ($id1, $id2, $cfg, $key)
}

rule Gh0st_RAT {
    meta:
        description  = "偵測 Gh0st RAT 網路封包標頭特徵"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://malpedia.caad.fkie.fraunhofer.de/details/win.gh0strat"

    strings:
        $header1 = "Gh0st" ascii
        $header2 = "HEART" ascii
        $cfg     = "gh0st.ini" ascii nocase

    condition:
        any of them
}
