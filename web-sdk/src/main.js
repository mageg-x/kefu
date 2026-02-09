import { createApp } from "vue";
import "./style.css";
import Widget from "./components/Widget.vue";

export function mountKefuWidget(props) {
  const div = document.createElement("div");
  div.id = "kefu-widget-root";
  div.style.position = "fixed";
  div.style.bottom = "5px";
  div.style.right = "5px";
  div.style.zIndex = "50";
  div.style.opacity = "0";
  div.style.transition = "opacity 0.3s ease";
  document.body.appendChild(div);

  const app = createApp(Widget, props);
  app.mount("#kefu-widget-root");

  setTimeout(() => {
    div.style.opacity = "1";
  }, 0);

  return app;
}

window.KefuChat = {
  init: (options = {}) => {
    console.log("KefuChat initializing with:", options);
    return mountKefuWidget(options);
  },
};

const scriptTags = document.querySelectorAll("script");
let appId = null;

for (const script of scriptTags) {
  if (script.hasAttribute("data-kefu-appid")) {
    appId = script.getAttribute("data-kefu-appid");
    break;
  }
}

if (appId) {
  window.KefuChat.init({ appId });
}
