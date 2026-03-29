rule FakeC2Beacon {
    meta:
        description = "Detects fake C2 beacon pattern for testing"
        severity = "high"
    strings:
        $c2 = "C2_SERVER=evil.example.com"
        $xor = "XOR_KEY=0xDEADBEEF"
    condition:
        $c2 and $xor
}

rule FakeMalwarePayload {
    meta:
        description = "Detects fake malware payload marker"
        severity = "medium"
    strings:
        $marker = "FAKE_MALWARE_PAYLOAD"
    condition:
        $marker
}
