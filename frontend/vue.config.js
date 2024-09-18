const { defineConfig } = require('@vue/cli-service');
module.exports = {
  devServer: {
    proxy: {
      // Proxy configuration for development
      '/api': {
        // Forward all requests starting with /api to the backend server
        target: 'http://localhost:8080',
        // Changes the origin of the host header to the target URL
        changeOrigin: true,
      },
    },
  },
};
