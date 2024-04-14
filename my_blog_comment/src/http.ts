import axios, { AxiosInstance } from 'axios';

class HttpClient {
  private instance: AxiosInstance | null = null;

  private getBaseURL(): string {
    if (process.env.NODE_ENV === 'development') {
      return 'http://localhost:8888/api';
    } else {
      const { protocol, hostname, port } = window.location;
      const portSegment = port && port !== '80' && port !== '443' ? `:${port}` : '';
      return `${protocol}//${hostname}${portSegment}/api`;
    }
  }

  public getInstance(): AxiosInstance {
    if (this.instance === null) {
      this.instance = axios.create({
        baseURL: this.getBaseURL(),
      });

      this.instance.interceptors.request.use();

      this.instance.interceptors.response.use();
    }

    return this.instance;
  }
}

export const httpClient = new HttpClient().getInstance();
