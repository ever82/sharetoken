#!/bin/bash
#
# Security Audit Script for ShareToken
# Runs automated security scans and generates reports
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SECURITY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$SECURITY_DIR")"

# Output directory
OUTPUT_DIR="$SECURITY_DIR/audits/$(date +%Y%m%d-%H%M%S)"
mkdir -p "$OUTPUT_DIR"

# Report files
REPORT_FILE="$OUTPUT_DIR/audit-report.md"
SUMMARY_FILE="$OUTPUT_DIR/summary.json"

echo "=========================================="
echo "ShareToken Security Audit"
echo "Started: $(date)"
echo "Output: $OUTPUT_DIR"
echo "=========================================="
echo

# Initialize report
cat > "$REPORT_FILE" << EOF
# Security Audit Report

**Date**: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
**Commit**: $(cd "$PROJECT_DIR" && git rev-parse --short HEAD 2>/dev/null || echo "N/A")
**Branch**: $(cd "$PROJECT_DIR" && git branch --show-current 2>/dev/null || echo "N/A")

## Summary

| Check | Status | Details |
|-------|--------|---------|
EOF

# Track results
declare -A RESULTS
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Function to add result to report
add_result() {
    local check_name="$1"
    local status="$2"
    local details="$3"

    echo "| $check_name | $status | $details |" >> "$REPORT_FILE"

    TOTAL_CHECKS=$((TOTAL_CHECKS + 1))
    if [ "$status" == "✅ PASS" ]; then
        PASSED_CHECKS=$((PASSED_CHECKS + 1))
    elif [ "$status" == "❌ FAIL" ]; then
        FAILED_CHECKS=$((FAILED_CHECKS + 1))
    fi
}

# Check 1: Go version
echo "[1/10] Checking Go version..."
GO_VERSION=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//' || echo "")
if [ -n "$GO_VERSION" ]; then
    MAJOR=$(echo "$GO_VERSION" | cut -d. -f1)
    MINOR=$(echo "$GO_VERSION" | cut -d. -f2)
    if [ "$MAJOR" -ge 1 ] && [ "$MINOR" -ge 19 ]; then
        echo -e "${GREEN}✓${NC} Go version: $GO_VERSION"
        add_result "Go Version" "✅ PASS" "$GO_VERSION"
    else
        echo -e "${RED}✗${NC} Go version too old: $GO_VERSION (need 1.19+)"
        add_result "Go Version" "❌ FAIL" "$GO_VERSION"
    fi
else
    echo -e "${RED}✗${NC} Go not found"
    add_result "Go Version" "❌ FAIL" "Not installed"
fi

# Check 2: Dependency vulnerabilities
echo "[2/10] Checking dependencies..."
cd "$PROJECT_DIR"
if command -v govulncheck &> /dev/null; then
    if govulncheck ./... > "$OUTPUT_DIR/govulncheck.log" 2>&1; then
        echo -e "${GREEN}✓${NC} No vulnerability findings"
        add_result "Dependency Vulns" "✅ PASS" "No issues found"
    else
        VULN_COUNT=$(grep -c "^Vulnerability" "$OUTPUT_DIR/govulncheck.log" 2>/dev/null || echo "0")
        echo -e "${YELLOW}!${NC} Found $VULN_COUNT vulnerabilities (see govulncheck.log)"
        add_result "Dependency Vulns" "⚠️ WARN" "$VULN_COUNT found"
    fi
else
    echo -e "${YELLOW}!${NC} govulncheck not installed"
    add_result "Dependency Vulns" "⚠️ SKIP" "Tool not installed"
fi

# Check 3: Go security scan (gosec)
echo "[3/10] Running gosec scan..."
if command -v gosec &> /dev/null; then
    if gosec -fmt json -out "$OUTPUT_DIR/gosec-report.json" ./... 2> "$OUTPUT_DIR/gosec.log"; then
        ISSUES=$(jq '.Stats.found // 0' "$OUTPUT_DIR/gosec-report.json" 2>/dev/null || echo "0")
        if [ "$ISSUES" -eq 0 ]; then
            echo -e "${GREEN}✓${NC} No security issues found"
            add_result "Go Security" "✅ PASS" "No issues"
        else
            echo -e "${YELLOW}!${NC} Found $ISSUES security issues"
            add_result "Go Security" "⚠️ WARN" "$ISSUES issues"
        fi
    else
        echo -e "${YELLOW}!${NC} gosec scan completed with findings"
        add_result "Go Security" "⚠️ WARN" "See gosec-report.json"
    fi
else
    echo -e "${YELLOW}!${NC} gosec not installed (go install github.com/securego/gosec/v2/cmd/gosec@latest)"
    add_result "Go Security" "⚠️ SKIP" "Tool not installed"
fi

# Check 4: Static analysis
echo "[4/10] Running static analysis..."
if command -v staticcheck &> /dev/null; then
    if staticcheck ./... > "$OUTPUT_DIR/staticcheck.log" 2>&1; then
        echo -e "${GREEN}✓${NC} Static analysis passed"
        add_result "Static Analysis" "✅ PASS" "No issues"
    else
        ISSUES=$(wc -l < "$OUTPUT_DIR/staticcheck.log" | tr -d ' ')
        echo -e "${YELLOW}!${NC} Found $ISSUES static analysis issues"
        add_result "Static Analysis" "⚠️ WARN" "$ISSUES issues"
    fi
else
    echo -e "${YELLOW}!${NC} staticcheck not installed"
    add_result "Static Analysis" "⚠️ SKIP" "Tool not installed"
fi

# Check 5: Test coverage
echo "[5/10] Checking test coverage..."
if go test -coverprofile="$OUTPUT_DIR/coverage.out" ./... > "$OUTPUT_DIR/test.log" 2>&1; then
    COVERAGE=$(go tool cover -func="$OUTPUT_DIR/coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
    COVERAGE_INT=${COVERAGE%.*}
    if [ "$COVERAGE_INT" -ge 70 ]; then
        echo -e "${GREEN}✓${NC} Test coverage: ${COVERAGE}%"
        add_result "Test Coverage" "✅ PASS" "${COVERAGE}%"
    else
        echo -e "${YELLOW}!${NC} Test coverage low: ${COVERAGE}% (target: 70%)"
        add_result "Test Coverage" "⚠️ WARN" "${COVERAGE}% (target 70%)"
    fi
else
    echo -e "${RED}✗${NC} Tests failed"
    add_result "Test Coverage" "❌ FAIL" "Tests failed"
fi

# Check 6: Secret scanning
echo "[6/10] Scanning for secrets..."
if command -v detect-secrets &> /dev/null; then
    if detect-secrets scan --all-files > "$OUTPUT_DIR/secrets.json" 2>/dev/null; then
        SECRET_COUNT=$(jq '.results | length' "$OUTPUT_DIR/secrets.json" 2>/dev/null || echo "0")
        if [ "$SECRET_COUNT" -eq 0 ]; then
            echo -e "${GREEN}✓${NC} No secrets found"
            add_result "Secret Scan" "✅ PASS" "No secrets"
        else
            echo -e "${YELLOW}!${NC} Found $SECRET_COUNT potential secrets"
            add_result "Secret Scan" "⚠️ WARN" "$SECRET_COUNT found"
        fi
    else
        add_result "Secret Scan" "⚠️ SKIP" "Scan failed"
    fi
else
    echo -e "${YELLOW}!${NC} detect-secrets not installed"
    add_result "Secret Scan" "⚠️ SKIP" "Tool not installed"
fi

# Check 7: File permissions
echo "[7/10] Checking file permissions..."
SENSITIVE_FILES=$(find "$PROJECT_DIR" -type f \( -name "*.key" -o -name "*.pem" -o -name "*private*" -o -name ".env*" \) ! -path "*/.git/*" 2>/dev/null)
PERM_ISSUES=0
for file in $SENSITIVE_FILES; do
    if [ -f "$file" ]; then
        PERM=$(stat -c "%a" "$file" 2>/dev/null || stat -f "%A" "$file" 2>/dev/null)
        if [ -n "$PERM" ] && [ "$PERM" -gt 600 ]; then
            echo -e "${YELLOW}!${NC} $file has permissions $PERM (should be 600)"
            PERM_ISSUES=$((PERM_ISSUES + 1))
        fi
    fi
done

if [ "$PERM_ISSUES" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} File permissions OK"
    add_result "File Permissions" "✅ PASS" "No issues"
else
    add_result "File Permissions" "⚠️ WARN" "$PERM_ISSUES files"
fi

# Check 8: Configuration security
echo "[8/10] Checking configuration files..."
CONFIG_ISSUES=0

# Check for hardcoded passwords
if grep -r "password.*=" --include="*.go" --include="*.yaml" --include="*.yml" "$PROJECT_DIR" 2>/dev/null | grep -v "// " | grep -v "# " | head -5 > "$OUTPUT_DIR/hardcoded.log"; then
    COUNT=$(wc -l < "$OUTPUT_DIR/hardcoded.log" | tr -d ' ')
    if [ "$COUNT" -gt 0 ]; then
        echo -e "${YELLOW}!${NC} Found $COUNT potential hardcoded passwords"
        CONFIG_ISSUES=$((CONFIG_ISSUES + COUNT))
    fi
fi

# Check for debug mode
if grep -r "debug.*true" --include="*.go" --include="*.yaml" "$PROJECT_DIR/config" 2>/dev/null > "$OUTPUT_DIR/debug-mode.log"; then
    echo -e "${YELLOW}!${NC} Debug mode enabled in configuration"
    CONFIG_ISSUES=$((CONFIG_ISSUES + 1))
fi

if [ "$CONFIG_ISSUES" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Configuration security OK"
    add_result "Configuration" "✅ PASS" "No issues"
else
    add_result "Configuration" "⚠️ WARN" "$CONFIG_ISSUES issues"
fi

# Check 9: Docker security
echo "[9/10] Checking Docker configurations..."
if [ -f "$PROJECT_DIR/Dockerfile" ]; then
    # Check for root user
    if grep -q "USER root" "$PROJECT_DIR/Dockerfile"; then
        echo -e "${YELLOW}!${NC} Dockerfile runs as root"
        add_result "Docker Security" "⚠️ WARN" "Runs as root"
    else
        echo -e "${GREEN}✓${NC} Dockerfile security OK"
        add_result "Docker Security" "✅ PASS" "No issues"
    fi
else
    add_result "Docker Security" "⚠️ SKIP" "No Dockerfile"
fi

# Check 10: Documentation
echo "[10/10] Checking security documentation..."
DOC_ISSUES=0

if [ ! -f "$SECURITY_DIR/README.md" ]; then
    echo -e "${YELLOW}!${NC} Security README missing"
    DOC_ISSUES=$((DOC_ISSUES + 1))
fi

if [ ! -f "$SECURITY_DIR/policies/security-policy.md" ]; then
    echo -e "${YELLOW}!${NC} Security policy missing"
    DOC_ISSUES=$((DOC_ISSUES + 1))
fi

if [ ! -f "$SECURITY_DIR/policies/incident-response.md" ]; then
    echo -e "${YELLOW}!${NC} Incident response plan missing"
    DOC_ISSUES=$((DOC_ISSUES + 1))
fi

if [ "$DOC_ISSUES" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Security documentation complete"
    add_result "Documentation" "✅ PASS" "Complete"
else
    add_result "Documentation" "⚠️ WARN" "$DOC_ISSUES missing"
fi

# Finalize report
cat >> "$REPORT_FILE" << EOF

## Detailed Results

### Dependency Vulnerabilities
\`\`\`
$(cat "$OUTPUT_DIR/govulncheck.log" 2>/dev/null || echo "No results")
\`\`\`

### Security Scan Results
See \`gosec-report.json\` for detailed findings.

### Test Output
\`\`\`
$(tail -50 "$OUTPUT_DIR/test.log" 2>/dev/null || echo "No test output")
\`\`\`

## Recommendations

EOF

if [ "$FAILED_CHECKS" -gt 0 ]; then
    echo "1. **Address failed checks immediately**" >> "$REPORT_FILE"
fi
if [ "$FAILED_CHECKS" -eq 0 ] && [ "$PASSED_CHECKS" -lt "$TOTAL_CHECKS" ]; then
    echo "1. **Review warnings and improve where possible**" >> "$REPORT_FILE"
fi
echo "2. **Install missing security tools** for more comprehensive scanning" >> "$REPORT_FILE"
echo "3. **Regularly update dependencies** to patch vulnerabilities" >> "$REPORT_FILE"
echo "4. **Review and update security policies** periodically" >> "$REPORT_FILE"

# Generate JSON summary
cat > "$SUMMARY_FILE" << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "total_checks": $TOTAL_CHECKS,
  "passed": $PASSED_CHECKS,
  "failed": $FAILED_CHECKS,
  "warnings": $((TOTAL_CHECKS - PASSED_CHECKS - FAILED_CHECKS)),
  "commit": "$(cd "$PROJECT_DIR" && git rev-parse HEAD 2>/dev/null || echo "N/A")",
  "results": {
EOF

# Add results to JSON
first=true
for key in "${!RESULTS[@]}"; do
    if [ "$first" = true ]; then
        first=false
    else
        echo "," >> "$SUMMARY_FILE"
    fi
    echo "    \"$key\": \"${RESULTS[$key]}\"" >> "$SUMMARY_FILE"
done

cat >> "$SUMMARY_FILE" << EOF

  }
}
EOF

# Print summary
echo
echo "=========================================="
echo "Audit Complete"
echo "=========================================="
echo -e "Total Checks:  $TOTAL_CHECKS"
echo -e "${GREEN}Passed:        $PASSED_CHECKS${NC}"
echo -e "${RED}Failed:        $FAILED_CHECKS${NC}"
echo -e "${YELLOW}Warnings:      $((TOTAL_CHECKS - PASSED_CHECKS - FAILED_CHECKS))${NC}"
echo
echo "Reports saved to: $OUTPUT_DIR"
echo "  - $REPORT_FILE"
echo "  - $SUMMARY_FILE"
echo "=========================================="

# Exit with error if any checks failed
if [ "$FAILED_CHECKS" -gt 0 ]; then
    echo -e "${RED}Audit completed with failures${NC}"
    exit 1
else
    echo -e "${GREEN}Audit completed successfully${NC}"
    exit 0
fi
