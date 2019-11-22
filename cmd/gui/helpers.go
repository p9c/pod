package gui

import (
	"fmt"
	"github.com/p9c/pod/pkg/log"
	"io/ioutil"
	"net/http"
)

func
getFile(f string, fs http.FileSystem) string {
	file, err := fs.Open(f)
	if err != nil {
		log.FATAL(err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	return string(body)
}

func evalJs(rc *rcvar) {
	var err error
	err = rc.WebView.Eval(getFile("libs/vue/vue.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("libs/vue/ej2-vue.min.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("libs/vue/vfg.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/duos.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/logo.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/overview.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/history.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/addressbook.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/explorer.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/ico/settings.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/balance.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/send.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/peers.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/status.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/networkhashrate.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/localhashrate.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/history.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/latestxs.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/addressbook.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/panels/settings.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/pages/overview.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/pages/history.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/pages/addressbook.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/pages/explorer.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/pages/settings.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/layout/header.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/layout/nav.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/layout/xorg.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = rc.WebView.Eval(getFile("w/js/dui.js",*rc.FileSystem))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}



}

func injectCss(rc *rcvar) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	rc.WebView.InjectCSS(getFile("w/css/material.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/theme/root.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/theme/colors.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/theme/grid.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/theme/helpers.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/theme/style.css",*rc.FileSystem))
	rc.WebView.InjectCSS(getFile("w/css/dui.css",*rc.FileSystem))

}