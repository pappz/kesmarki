package com.kesmarki;

import android.content.Context;
import android.webkit.JavascriptInterface;

public class WebViewInterface {
   Context context;

   WebViewInterface(Context context) {
      this.context = context;
   }

   @JavascriptInterface
   public String getPassword() {
      return BuildConfig.MQTT_PASSWORD;
   }
}