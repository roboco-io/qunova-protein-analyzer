#!/bin/bash

# API Test Script for Protein Analyzer Backend
# Usage: ./test_api.sh

API_URL="http://localhost:8080/api/v1"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Testing Protein Analyzer API..."
echo "================================"

# Function to print colored output
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
    fi
}

# 1. Health Check
echo -e "\n${YELLOW}1. Testing Health Check${NC}"
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
if [ "$response" = "200" ]; then
    print_result 0 "Health check passed"
else
    print_result 1 "Health check failed (HTTP $response)"
    echo "Make sure the server is running on port 8080"
    exit 1
fi

# 2. File Upload (using a sample PDB file or create a dummy one)
echo -e "\n${YELLOW}2. Testing File Upload${NC}"

# Create a minimal PDB file for testing
cat > test_protein.pdb << 'EOF'
HEADER    TEST PROTEIN                            01-JAN-24   TEST
ATOM      1  N   ALA A   1      10.000  10.000  10.000  1.00 20.00           N
ATOM      2  CA  ALA A   1      11.000  11.000  11.000  1.00 20.00           C
ATOM      3  C   ALA A   1      12.000  12.000  12.000  1.00 20.00           C
ATOM      4  O   ALA A   1      13.000  13.000  13.000  1.00 20.00           O
END
EOF

upload_response=$(curl -s -X POST \
    -F "file=@test_protein.pdb" \
    ${API_URL}/upload)

file_id=$(echo $upload_response | grep -o '"fileId":"[^"]*' | sed 's/"fileId":"//')

if [ ! -z "$file_id" ]; then
    print_result 0 "File uploaded successfully (ID: $file_id)"
else
    print_result 1 "File upload failed"
    echo "Response: $upload_response"
fi

# 3. Start Analysis
echo -e "\n${YELLOW}3. Testing Analysis Start${NC}"
analysis_response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"options":{"includeSequence":true,"calculateDomains":true,"identifyLigands":true}}' \
    ${API_URL}/analyze/${file_id})

analysis_id=$(echo $analysis_response | grep -o '"analysisId":"[^"]*' | sed 's/"analysisId":"//')

if [ ! -z "$analysis_id" ]; then
    print_result 0 "Analysis started (ID: $analysis_id)"
else
    print_result 1 "Analysis start failed"
    echo "Response: $analysis_response"
fi

# 4. Check Analysis Status
echo -e "\n${YELLOW}4. Testing Analysis Status${NC}"
sleep 2  # Wait for analysis to complete

status_response=$(curl -s ${API_URL}/analysis/${analysis_id})
status=$(echo $status_response | grep -o '"status":"[^"]*' | sed 's/"status":"//')

if [ "$status" = "completed" ] || [ "$status" = "processing" ]; then
    print_result 0 "Analysis status check successful (Status: $status)"
else
    print_result 1 "Analysis status check failed"
    echo "Response: $status_response"
fi

# 5. Get Chains Data
echo -e "\n${YELLOW}5. Testing Chains Retrieval${NC}"
chains_response=$(curl -s ${API_URL}/analysis/${analysis_id}/chains)

if echo "$chains_response" | grep -q "chains"; then
    print_result 0 "Chains data retrieved successfully"
else
    print_result 1 "Chains retrieval failed"
fi

# 6. Get Domains Data
echo -e "\n${YELLOW}6. Testing Domains Retrieval${NC}"
domains_response=$(curl -s ${API_URL}/analysis/${analysis_id}/domains)

if echo "$domains_response" | grep -q "domains"; then
    print_result 0 "Domains data retrieved successfully"
else
    print_result 1 "Domains retrieval failed"
fi

# 7. Get Ligands Data
echo -e "\n${YELLOW}7. Testing Ligands Retrieval${NC}"
ligands_response=$(curl -s ${API_URL}/analysis/${analysis_id}/ligands)

if echo "$ligands_response" | grep -q "ligands"; then
    print_result 0 "Ligands data retrieved successfully"
else
    print_result 1 "Ligands retrieval failed"
fi

# 8. Export Data
echo -e "\n${YELLOW}8. Testing Data Export${NC}"
export_response=$(curl -s -o exported_data.json -w "%{http_code}" \
    ${API_URL}/export/${analysis_id})

if [ "$export_response" = "200" ]; then
    print_result 0 "Data exported successfully (saved to exported_data.json)"
else
    print_result 1 "Data export failed (HTTP $export_response)"
fi

# Cleanup
rm -f test_protein.pdb exported_data.json

echo -e "\n================================"
echo "✅ API testing completed!"
echo ""
echo "To test with a real PDB file:"
echo "curl -X POST -F 'file=@your_protein.pdb' ${API_URL}/upload"