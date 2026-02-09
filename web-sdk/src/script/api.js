import axios from "axios";

class Api {
  constructor() {
    this.baseURL = "http://localhost:5300";
    this.userId = null;
    this.api = axios.create({
      baseURL: this.baseURL,
      timeout: 10000,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  setBaseURL(url) {
    this.baseURL = url;
    this.api.defaults.baseURL = url;
  }

  setUserId(userId) {
    this.userId = userId;
  }

  async getConfig(appId) {
    try {
      const params = {
        appid: appId,
      };
      
      if (this.userId) {
        params.userid = this.userId;
      }
      
      const response = await this.api.get("/api/v1/config", {
        params: params,
      });
      return response.data;
    } catch (error) {
      console.error("Failed to get config:", error);
      throw new Error(error.response?.data?.msg || "获取配置失败");
    }
  }
}

export default new Api();