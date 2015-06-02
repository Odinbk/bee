package utils

import (
    "bitbucket.org/tebeka/selenium"
)

func Login(homeUrl, username, password string) (session_id string) {
    caps := selenium.Capabilities{"browserName": "firefox"}
    wd, _ := selenium.NewRemote(caps, "")
    defer wd.Quit()

    usernameIptSelector := "input#username"
    passwordIptSelector := "input#password"
    loginBtnSelector := "#login_submit_button"

    wd.Get(homeUrl)
    input(wd, usernameIptSelector, username)
    input(wd, passwordIptSelector, password)
    loginBtn, _ := wd.FindElement(selenium.ByCSSSelector, loginBtnSelector)
    loginBtn.Click()

    cookie, _ := wd.GetCookies()
    for _, item := range(cookie) {
        if item.Name == "sessionId" {
            session_id = item.Value
            break
        }
    }
    return session_id
}

func input(wd selenium.WebDriver,  selector, content string) {
    elem, _ := wd.FindElement(selenium.ByCSSSelector, selector)
    elem.Clear()
    elem.SendKeys(content)
}
