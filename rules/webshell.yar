/*
    網頁後門（Webshell）特徵規則集
    涵蓋：常見 PHP、ASP、JSP webshell 特徵
*/

rule PHP_Webshell_Generic {
    meta:
        description  = "偵測常見 PHP webshell：eval + base64 或 system/exec 組合"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://github.com/xl7dev/WebShell"

    strings:
        $eval    = "eval(" ascii nocase
        $b64     = "base64_decode(" ascii nocase
        $system  = "system(" ascii nocase
        $exec    = "exec(" ascii nocase
        $shell   = "shell_exec(" ascii nocase
        $passthru = "passthru(" ascii nocase
        $assert  = "assert(" ascii nocase
        $input   = "$_REQUEST" ascii nocase
        $input2  = "$_POST" ascii nocase
        $input3  = "$_GET" ascii nocase

    condition:
        ($eval and $b64) or
        (($system or $exec or $shell or $passthru or $assert) and ($input or $input2 or $input3))
}

rule PHP_Webshell_China_Chopper {
    meta:
        description  = "偵測 China Chopper webshell 特徵（eval+assert 單行結構）"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://www.cisa.gov/news-events/alerts/2019/09/06/malware-analysis-report-chinachopper"

    strings:
        $cc1 = "eval(base64_decode($_POST" ascii nocase
        $cc2 = "assert($_POST" ascii nocase
        $cc3 = "eval(gzinflate(base64_decode(" ascii nocase

    condition:
        any of them
}

rule ASP_Webshell_Generic {
    meta:
        description  = "偵測常見 ASP/ASPX webshell：CreateObject + WScript.Shell 組合"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1505/003/"

    strings:
        $obj1 = "CreateObject(\"WScript.Shell\")" ascii wide nocase
        $obj2 = "CreateObject(\"Scripting.FileSystemObject\")" ascii wide nocase
        $exec = ".Run(" ascii wide nocase
        $req  = "Request(" ascii wide nocase
        $req2 = "Request.Form(" ascii wide nocase

    condition:
        ($obj1 or $obj2) and ($exec) and ($req or $req2)
}

rule JSP_Webshell_Runtime_Exec {
    meta:
        description  = "偵測 JSP webshell 透過 Runtime.exec() 執行系統命令"
        severity     = "high"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1505/003/"

    strings:
        $rt  = "Runtime.getRuntime().exec(" ascii nocase
        $req = "request.getParameter(" ascii nocase
        $cmd = "cmd" ascii nocase

    condition:
        $rt and $req and $cmd
}
