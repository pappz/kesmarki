package com.kesmarki;

import android.content.Context;
import android.webkit.JavascriptInterface;

public class WebViewInterface {
   Context context;

   WebViewInterface(Context c) {
      context = c;
   }

   @JavascriptInterface
   public String getPassword() {
      return "secret password";
   }
}