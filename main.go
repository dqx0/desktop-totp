/**
 *  Copyright 2024 dqx0
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/getlantern/systray"
	"github.com/pquerna/otp/totp"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icoToByteArray(os.Getenv("TOTP_ICON_PATH")))
	systray.SetTitle("TOTP")
	systray.SetTooltip("TOTP")

	secret := os.Getenv("TOTP_SECRET")

	if secret == "" {
		systray.AddMenuItem("No secret found", "No secret found")
		return
	}

	copyMenuItem := systray.AddMenuItem("Copy", "Copy the TOTP code to the clipboard")
	exitMenuItem := systray.AddMenuItem("Exit", "Exit the application")

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			code, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				fmt.Println("Error generating TOTP code: ", err)
				continue
			}
			systray.SetTooltip(code)
		}
	}()

	go func() {
		for range copyMenuItem.ClickedCh {
			otp, err := totp.GenerateCode(secret, time.Now())
			if err != nil {
				fmt.Println("Error generating TOTP code: ", err)
				continue
			}
			err = clipboard.WriteAll(otp)
			if err != nil {
				fmt.Println("Error copying TOTP code to clipboard: ", err)
			} else {
				fmt.Println("Copied TOTP code to clipboard")
			}
		}
	}()
	<-exitMenuItem.ClickedCh
	systray.Quit()
}

func onExit() {
}

func icoToByteArray(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	byteArray, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	return byteArray
}
