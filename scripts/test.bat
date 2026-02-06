@echo off
chcp 936 >nul
echo ============================================
echo Geelato CLI Test Script
echo ============================================
echo.

echo [Test 1] Show help
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe --help
echo.
pause

echo [Test 2] Show version
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe --version
echo.
pause

echo [Test 3] Config set
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe config set test.key "hello world"
echo.

echo [Test 4] Config get
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe config get test.key
echo.

echo [Test 5] Config list
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe config list
echo.
pause

echo [Test 6] Verbose mode
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe --verbose app
echo.
pause

echo [Test 7] JSON mode
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe --json config list
echo.
pause

echo [Test 8] Config remove
d:\geelato\geelato-doc\geelato-cli\bin\geelato.exe config remove test.key
echo.

echo ============================================
echo All tests completed!
echo ============================================
