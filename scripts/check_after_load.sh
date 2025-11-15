$ErrorActionPreference = "Stop"
$baseUrl = "http://localhost:8080"

Write-Host "=== 1. Checking /health ==="
$health = Invoke-WebRequest -Uri "$baseUrl/health" -Method GET
Write-Host "health status:" $health.StatusCode
if ($health.StatusCode -ne 200) {
    Write-Error "ERROR: /health returned $($health.StatusCode)"
    exit 1
}

Write-Host ""
Write-Host "=== 2. Checking /stats/assignments ==="
$stats = Invoke-WebRequest -Uri "$baseUrl/stats/assignments" -Method GET
$stats.Content | Write-Host

Write-Host ""
Write-Host "=== 3. Creating 10 PRs (author_id = u1..u5) and merging them ==="

$success = 0
$failed = 0

for ($i = 1; $i -le 10; $i++) {
    $prId = "postcheck-pr-$i"
    $authorIndex = ($i % 5) + 1
    $authorId = "u$authorIndex"

    Write-Host "--- PR $prId (author: $authorId) ---"

    $createBody = @{
        pull_request_id   = $prId
        pull_request_name = "Post-check PR $i"
        author_id         = $authorId
    } | ConvertTo-Json

    try {
        $createResp = Invoke-WebRequest -Uri "$baseUrl/pullRequest/create" `
            -Method POST `
            -ContentType "application/json" `
            -Body $createBody `
            -ErrorAction Stop
        $createCode = $createResp.StatusCode
    } catch {
        $createCode = $_.Exception.Response.StatusCode.Value__
    }

    Write-Host "create status:" $createCode

    if ($createCode -eq 201 -or $createCode -eq 409) {
        $mergeBody = @{ pull_request_id = $prId } | ConvertTo-Json
        try {
            $mergeResp = Invoke-WebRequest -Uri "$baseUrl/pullRequest/merge" `
                -Method POST `
                -ContentType "application/json" `
                -Body $mergeBody `
                -ErrorAction Stop
            $mergeCode = $mergeResp.StatusCode
        } catch {
            $mergeCode = $_.Exception.Response.StatusCode.Value__
        }

        Write-Host "merge status:" $mergeCode

        if ($mergeCode -eq 200 -or $mergeCode -eq 404) {
            $success++
        } else {
            Write-Host "ERROR: unexpected merge status: $mergeCode"
            $failed++
        }
    } else {
        Write-Host "ERROR: unexpected create status: $createCode"
        $failed++
    }
}

Write-Host ""
Write-Host "Create+merge scenarios: success=$success, failed=$failed"
if ($failed -gt 0) {
    Write-Error "ERROR: some post-check scenarios failed"
    exit 1
}

Write-Host ""
Write-Host "=== 4. Checking /users/getReview for several users ==="

foreach ($uid in @("u1", "u2", "u3", "u4", "u5")) {
    Write-Host "--- user_id=$uid ---"
    try {
        $resp = Invoke-WebRequest -Uri "$baseUrl/users/getReview?user_id=$uid" `
            -Method GET `
            -ErrorAction Stop
        $code = $resp.StatusCode
        $body = $resp.Content
    } catch {
        $code = $_.Exception.Response.StatusCode.Value__
        $body = $_.Exception.Response.StatusDescription
    }

    Write-Host "status:" $code
    Write-Host "body:  " $body

    if ($code -ne 200) {
        Write-Error "ERROR: /users/getReview for $uid returned $code"
        exit 1
    }
}

Write-Host ""
Write-Host "=== 5. Final /stats/assignments ==="
$stats2 = Invoke-WebRequest -Uri "$baseUrl/stats/assignments" -Method GET
$stats2.Content | Write-Host

Write-Host "=== POST-CHECK OK ==="
