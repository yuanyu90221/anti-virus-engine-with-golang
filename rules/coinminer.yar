/*
    加密貨幣挖礦程式（Coinminer）特徵規則集
    涵蓋：XMRig、常見挖礦池連線字串、CPU 竊用行為特徵
*/

rule XMRig_Miner {
    meta:
        description  = "偵測 XMRig Monero 挖礦工具特徵字串"
        severity     = "medium"
        date         = "2026-03-29"
        reference    = "https://github.com/xmrig/xmrig"

    strings:
        $id1  = "xmrig" ascii wide nocase
        $id2  = "XMRig" ascii
        $pool = "stratum+tcp://" ascii wide nocase
        $pool2 = "stratum+ssl://" ascii wide nocase
        $cfg  = "\"donate-level\"" ascii wide nocase
        $cfg2 = "\"coin\": \"XMR\"" ascii wide nocase

    condition:
        $id1 or $id2 or ($pool and ($cfg or $cfg2)) or ($pool2 and ($cfg or $cfg2))
}

rule Generic_Coinminer_Pool {
    meta:
        description  = "偵測連接常見 Monero 挖礦池的字串特徵"
        severity     = "medium"
        date         = "2026-03-29"
        reference    = "https://attack.mitre.org/techniques/T1496/"

    strings:
        $pool1 = "pool.minexmr.com" ascii wide nocase
        $pool2 = "xmrpool.eu" ascii wide nocase
        $pool3 = "monerohash.com" ascii wide nocase
        $pool4 = "c3pool.com" ascii wide nocase
        $pool5 = "supportxmr.com" ascii wide nocase
        $stratum = "stratum+tcp://" ascii wide nocase

    condition:
        any of ($pool1, $pool2, $pool3, $pool4, $pool5) or $stratum
}

rule Cryptojacking_Script {
    meta:
        description  = "偵測瀏覽器端 JavaScript 挖礦腳本（CoinHive 衍生）"
        severity     = "medium"
        date         = "2026-03-29"
        reference    = "https://www.malwarebytes.com/cryptojacking"

    strings:
        $ch1 = "CoinHive" ascii nocase
        $ch2 = "coinhive.min.js" ascii nocase
        $wasm = "cryptonight.wasm" ascii nocase
        $worker = "CryptoNight" ascii nocase

    condition:
        any of them
}
