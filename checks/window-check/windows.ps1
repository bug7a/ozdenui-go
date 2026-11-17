# The real window of OzdenUI on Windows (the web view and the browser window), run by a scheduled task in the session
# of the user: a program started over SSH has no desktop. probe.exe is built next to this file:
#   go build -o checks/window-check/probe.exe ./checks/window-check
# It writes probe.log and a picture of each window here.
$ErrorActionPreference = "Continue"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$log = Join-Path $here "probe.log"
"start $(Get-Date -Format s)" | Set-Content $log
Add-Type -AssemblyName System.Drawing
Add-Type @"
using System; using System.Runtime.InteropServices;
public class W {
  [DllImport("user32.dll")] public static extern bool SetProcessDPIAware();
  [DllImport("user32.dll")] public static extern bool GetWindowRect(IntPtr h, out R r);
  [DllImport("user32.dll")] public static extern bool GetClientRect(IntPtr h, out R r);
  [DllImport("user32.dll")] public static extern uint GetDpiForWindow(IntPtr h);
  [DllImport("user32.dll")] public static extern bool SetForegroundWindow(IntPtr h);
  [DllImport("user32.dll")] public static extern IntPtr FindWindow(string c, string t);
  public struct R { public int L, T, Rt, B; }
}
"@
[W]::SetProcessDPIAware() | Out-Null

function Shot($hwnd, $name) {
  [W]::SetForegroundWindow($hwnd) | Out-Null
  Start-Sleep -Milliseconds 500
  $r = New-Object W+R; [W]::GetWindowRect($hwnd, [ref]$r) | Out-Null
  $c = New-Object W+R; [W]::GetClientRect($hwnd, [ref]$c) | Out-Null
  "  $name window $($r.Rt-$r.L) x $($r.B-$r.T) at $($r.L),$($r.T); inside $($c.Rt) x $($c.B); dpi $([W]::GetDpiForWindow($hwnd))" | Add-Content $log
  $bmp = New-Object System.Drawing.Bitmap ($r.Rt-$r.L), ($r.B-$r.T)
  $g = [System.Drawing.Graphics]::FromImage($bmp)
  $g.CopyFromScreen($r.L, $r.T, 0, 0, $bmp.Size)
  $bmp.Save((Join-Path $here "$name.png")); $g.Dispose(); $bmp.Dispose()
}

function WaitTitle($title, $seconds) {
  for ($i = 0; $i -lt $seconds * 4; $i++) {
    $w = Get-Process probe, msedge, chrome -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowTitle -like "Probe*" } | Select-Object -First 1
    if ($w -and $w.MainWindowTitle -eq $title) { "  title: '$($w.MainWindowTitle)' ($($w.ProcessName))" | Add-Content $log; return $w.MainWindowHandle }
    if ($w) { $script:lastTitle = "$($w.MainWindowTitle) ($($w.ProcessName))" }
    Start-Sleep -Milliseconds 250
  }
  "  no window titled '$title'; the last title seen: '$script:lastTitle'" | Add-Content $log
  return [IntPtr]::Zero
}

$exe = Join-Path $here "probe.exe"

# 1. The web view (WebView2), closed by the program: ui.Close() after 5 s.
$env:OZDENUI_WINDOW = ""; $env:PROBE_CLOSE = "1"; $env:PROBE_REPORT = Join-Path $here "page-webview.txt"
$p = Start-Process $exe -PassThru
$h = WaitTitle "Probe 2" 15
"1 webview: title 'Probe 2' found: $($h -ne [IntPtr]::Zero)" | Add-Content $log
if ($h -ne [IntPtr]::Zero) { Start-Sleep 1; Shot $h "webview" }
$ended = $p.WaitForExit(12000)
"1 webview: the program ended after Close(): $ended" | Add-Content $log
if (-not $ended) { $p.Kill() }
"  page: $(Get-Content $env:PROBE_REPORT -ErrorAction SilentlyContinue)" | Add-Content $log

# 2. The web view, closed by the user (the close button of the window: WM_CLOSE).
$env:PROBE_CLOSE = ""
$p = Start-Process $exe -PassThru
$h = WaitTitle "Probe 2" 15
Start-Sleep 1
$p.CloseMainWindow() | Out-Null
$ended = $p.WaitForExit(8000)
"2 webview: the program ended when its window was closed: $ended" | Add-Content $log
if (-not $ended) { $p.Kill() }

# 3. The browser window (OZDENUI_WINDOW=browser: what a computer without WebView2 gets), closed by the program.
$env:OZDENUI_WINDOW = "browser"; $env:PROBE_CLOSE = "1"; $env:PROBE_REPORT = Join-Path $here "page-browser.txt"
$p = Start-Process $exe -PassThru
$h = WaitTitle "Probe 2" 20
"3 browser: title 'Probe 2' found: $($h -ne [IntPtr]::Zero)" | Add-Content $log
if ($h -ne [IntPtr]::Zero) { Start-Sleep 1; Shot $h "browser" }
$ended = $p.WaitForExit(12000)
"3 browser: the program ended after Close(): $ended" | Add-Content $log
if (-not $ended) { $p.Kill() }
"  page: $(Get-Content $env:PROBE_REPORT -ErrorAction SilentlyContinue)" | Add-Content $log
"done $(Get-Date -Format s)" | Add-Content $log
