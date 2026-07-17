package com.kesmarki;

import androidx.appcompat.app.AppCompatActivity;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.net.Uri;
import android.os.Bundle;
import android.util.Log;
import android.webkit.ConsoleMessage;
import android.webkit.ValueCallback;
import android.webkit.WebChromeClient;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import androidx.core.splashscreen.SplashScreen;

public class MainActivity extends AppCompatActivity {

    private WebView webView;
    private final int FILECHOOSER_RESULTCODE = 56;
    private ValueCallback<Uri[]> filePathCallback;

    // Keep the launcher-icon splash on screen until the WebView finishes loading.
    private boolean webViewLoaded = false;

    @SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        SplashScreen splashScreen = SplashScreen.installSplashScreen(this);
        splashScreen.setKeepOnScreenCondition(() -> !webViewLoaded);

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        WebView.setWebContentsDebuggingEnabled(true);

        webView = findViewById(R.id.activity_main_webview);

        WebSettings webSettings = webView.getSettings();
        webSettings.setJavaScriptEnabled(true);
        webSettings.setDomStorageEnabled(true);
        webSettings.setAllowContentAccess(true);
        webSettings.setAllowFileAccess(true);
        webSettings.setMixedContentMode(WebSettings.MIXED_CONTENT_COMPATIBILITY_MODE);

        webView.clearCache(true);

        webView.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                view.loadUrl(request.getUrl().toString());
                return false;
            }

            @Override
            public void onPageFinished(WebView view, String url) {
                super.onPageFinished(view, url);
                // Page is rendered: let the splash screen dismiss.
                webViewLoaded = true;
            }
        });

        webView.setWebChromeClient(new WebChromeClient() {
            public boolean onShowFileChooser(WebView webView, ValueCallback<Uri[]> filePathCallback, WebChromeClient.FileChooserParams fileChooserParams) {
                MainActivity.this.filePathCallback = filePathCallback;
                MainActivity.this.startActivityForResult(fileChooserParams.createIntent(), FILECHOOSER_RESULTCODE);
                return true;
            }
            @Override
            public boolean onConsoleMessage(ConsoleMessage consoleMessage) {
                Log.d("Kesmarki-webview", consoleMessage.message() + " (" +
                        consoleMessage.lineNumber() + "): " + consoleMessage.sourceId());
                return true;
            }
        });

        setUserAgent();
        webView.addJavascriptInterface(new WebViewInterface(this), "KesmarkiApp");
        webView.getSettings().setJavaScriptEnabled(true);

        webView.loadUrl("https://kesmarki-46126.web.app");
    }

    @Override
    public void onActivityResult(int requestCode, int resultCode, Intent intent) {
        super.onActivityResult(requestCode, resultCode, intent);
        if (requestCode != FILECHOOSER_RESULTCODE) {
           return;
        }
        if(filePathCallback == null) {
            return;
        }
        if(resultCode != RESULT_OK) {
            return;
        }
        if(intent == null) {
            return;
        }

        Uri[] uris = WebChromeClient.FileChooserParams.parseResult(resultCode, intent);
        filePathCallback.onReceiveValue(uris);
        filePathCallback = null;
    }

    private void setUserAgent() {
        WebSettings webSettings = webView.getSettings();
        String userAgent = String.format("%s [%s %s]", webSettings.getUserAgentString(),
                "KesmarkiApp", BuildConfig.VERSION_CODE);
        webSettings.setUserAgentString(userAgent);
    }

}