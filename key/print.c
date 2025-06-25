#include <windows.h>
#include <stdio.h>
#include "print.h"

HHOOK hHook = NULL;

void VkCodeToUtf8(DWORD vkCode, char* outChar, int outSize) {
    BYTE keyboardState[256];
    if (!GetKeyboardState(keyboardState)) {
        snprintf(outChar, outSize, "?");
        return;
    }

    HWND foreground = GetForegroundWindow();
    DWORD threadId = GetWindowThreadProcessId(foreground, NULL);
    HKL layout = GetKeyboardLayout(threadId);

    UINT scanCode = MapVirtualKeyEx(vkCode, MAPVK_VK_TO_VSC, layout);

    WCHAR unicodeChar[5] = {0};
    int len = ToUnicodeEx(vkCode, scanCode, keyboardState, unicodeChar, 4, 0, layout);

    if (len > 0) {
        int bytesWritten = WideCharToMultiByte(CP_UTF8, 0, unicodeChar, len, outChar, outSize - 1, NULL, NULL);
        if (bytesWritten > 0) {
            outChar[bytesWritten] = '\0';
        } else {
            snprintf(outChar, outSize, "?");
        }
    } else {
        snprintf(outChar, outSize, "");
    }
}

LRESULT CALLBACK KeyboardProc(int nCode, WPARAM wParam, LPARAM lParam) {
    if (nCode >= 0) {
        KBDLLHOOKSTRUCT* p = (KBDLLHOOKSTRUCT*)lParam;
        if (wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN) {
            char text[16];
            VkCodeToUtf8(p->vkCode, text, sizeof(text));
            printf("Pressed key char: %s\n", text);
            GoHandleKey(text);
        }
    }
    return CallNextHookEx(hHook, nCode, wParam, lParam);
}

int printKeyBoard() {
    hHook = SetWindowsHookEx(WH_KEYBOARD_LL, KeyboardProc, NULL, 0);

    if (hHook == NULL) {
        printf("Error by connection hook\n");
        return 1;
    }

    MSG msg;
    while (GetMessage(&msg, NULL, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessage(&msg);
    }

    UnhookWindowsHookEx(hHook);
    return 0;
}
