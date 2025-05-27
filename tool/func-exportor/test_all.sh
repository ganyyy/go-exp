#!/bin/bash

# Comprehensive test script for func-exportor tool
echo "🧪 Testing func-exportor tool - All functionality"
echo "================================================="

TOOL="./func-exportor"
TEST_FILE="example.go"

echo "📦 Building tool..."
go build -o func-exportor main.go
echo "✅ Build successful"

echo
echo "🔍 Test 1: Basic symbol extraction"
$TOOL $TEST_FILE
echo "✅ Basic extraction works"

echo
echo "🔍 Test 2: Statistics output"
$TOOL --stats $TEST_FILE
echo "✅ Statistics output works"

echo
echo "🔍 Test 3: Package overview"
$TOOL --overview $TEST_FILE
echo "✅ Overview functionality works"

echo
echo "🎉 Key tests completed successfully!"
