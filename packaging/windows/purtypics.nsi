; NSIS Script for Purtypics
; Can be compiled on Linux/macOS using makensis

!define PRODUCT_NAME "Purtypics"
!define PRODUCT_VERSION "${VERSION}"
!define PRODUCT_PUBLISHER "Purtypics Contributors"
!define PRODUCT_WEB_SITE "https://github.com/cjs/purtypics"
!define PRODUCT_UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${PRODUCT_NAME}"
!define PRODUCT_UNINST_ROOT_KEY "HKLM"

SetCompressor lzma

; MUI Settings
!include "MUI2.nsh"
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; Welcome page
!insertmacro MUI_PAGE_WELCOME
; License page (skip if no LICENSE file)
; !insertmacro MUI_PAGE_LICENSE "..\..\LICENSE"
; Directory page
!insertmacro MUI_PAGE_DIRECTORY
; Instfiles page
!insertmacro MUI_PAGE_INSTFILES
; Finish page
!define MUI_FINISHPAGE_RUN "$INSTDIR\purtypics.exe"
!define MUI_FINISHPAGE_RUN_PARAMETERS "--version"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_INSTFILES

; Language files
!insertmacro MUI_LANGUAGE "English"

; Installer attributes
Name "${PRODUCT_NAME} ${PRODUCT_VERSION}"
OutFile "purtypics-installer.exe"
InstallDir "$PROGRAMFILES64\Purtypics"
InstallDirRegKey HKLM "${PRODUCT_UNINST_KEY}" "UninstallString"
ShowInstDetails show
ShowUnInstDetails show
RequestExecutionLevel admin

Section "MainSection" SEC01
  SetOutPath "$INSTDIR"
  SetOverwrite ifnewer
  
  ; Copy files
  File /oname=purtypics.exe "purtypics.exe"
  File "..\..\README.md"
  ; File /nonfatal "..\..\LICENSE"
  
  ; Create shortcuts
  CreateDirectory "$SMPROGRAMS\Purtypics"
  CreateShortcut "$SMPROGRAMS\Purtypics\Purtypics.lnk" "$INSTDIR\purtypics.exe"
  CreateShortcut "$DESKTOP\Purtypics.lnk" "$INSTDIR\purtypics.exe"
SectionEnd

Section -AdditionalIcons
  CreateShortcut "$SMPROGRAMS\Purtypics\Uninstall.lnk" "$INSTDIR\uninst.exe"
SectionEnd

Section -Post
  WriteUninstaller "$INSTDIR\uninst.exe"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayName" "$(^Name)"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "UninstallString" "$INSTDIR\uninst.exe"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayIcon" "$INSTDIR\purtypics.exe"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "DisplayVersion" "${PRODUCT_VERSION}"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "URLInfoAbout" "${PRODUCT_WEB_SITE}"
  WriteRegStr HKLM "${PRODUCT_UNINST_KEY}" "Publisher" "${PRODUCT_PUBLISHER}"
  
  ; Add to PATH
  Push "$INSTDIR"
  Call AddToPath
SectionEnd

; Functions
Function AddToPath
  Exch $0
  Push $1
  Push $2
  Push $3

  ReadRegStr $1 HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "PATH"
  StrCpy $2 "$1;$0"
  WriteRegExpandStr HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "PATH" "$2"
  SendMessage ${HWND_BROADCAST} ${WM_WININICHANGE} 0 "STR:Environment" /TIMEOUT=5000

  Pop $3
  Pop $2
  Pop $1
  Pop $0
FunctionEnd

Section Uninstall
  Delete "$INSTDIR\uninst.exe"
  Delete "$INSTDIR\purtypics.exe"
  Delete "$INSTDIR\README.md"
  ; Delete "$INSTDIR\LICENSE"

  Delete "$SMPROGRAMS\Purtypics\Uninstall.lnk"
  Delete "$SMPROGRAMS\Purtypics\Purtypics.lnk"
  Delete "$DESKTOP\Purtypics.lnk"

  RMDir "$SMPROGRAMS\Purtypics"
  RMDir "$INSTDIR"

  DeleteRegKey HKLM "${PRODUCT_UNINST_KEY}"
  SetAutoClose true
SectionEnd