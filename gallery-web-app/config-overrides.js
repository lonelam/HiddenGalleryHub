const { override, addLessLoader, adjustStyleLoaders, addWebpackResolve, removeModuleScopePlugin } = require("customize-cra");

module.exports = {
  // The Webpack config to use when compiling your react app for development or production.
  webpack: override(
    addWebpackResolve({
      fallback: {
        "path": require.resolve("path-browserify")
      }
    }),
    removeModuleScopePlugin(),
    addLessLoader({
    lessOptions: {
      strictMath: true,
      noIeCompat: true,
      modifyVars: {
        "@primary-color": "#1DA57A", // for example, you use Ant Design to change theme color.
      },
      cssLoaderOptions: {}, // .less file used css-loader option, not all CSS file.
      cssModules: {
        localIdentName: "[path][name]__[local]--[hash:base64:5]", // if you use CSS Modules, and custom `localIdentName`, default is '[local]--[hash:base64:5]'.
      }
    },
  }),
    adjustStyleLoaders(({ use: [, , postcss] }) => {
      const postcssOptions = postcss.options;
      postcss.options = { postcssOptions };
    })),
  // The function to use to create a webpack dev server configuration when running the development
  // server with 'npm run start' or 'yarn start'.
  // Example: set the dev server to use a specific certificate in https.
  devServer: function (configFunction) {
    // Return the replacement function for create-react-app to use to generate the Webpack
    // Development Server config. "configFunction" is the function that would normally have
    // been used to generate the Webpack Development server config - you can use it to create
    // a starting configuration to then modify instead of having to create a config from scratch.
    return function (proxy, allowedHost) {
      // Create the default config by calling configFunction with the proxy/allowedHost parameters
      const config = configFunction(proxy, allowedHost);

      // Change the https certificate options to match your certificate, using the .env file to
      // set the file paths & passphrase.
      // const fs = require('fs');
      // config.https = {
      //   key: fs.readFileSync(process.env.REACT_HTTPS_KEY, 'utf8'),
      //   cert: fs.readFileSync(process.env.REACT_HTTPS_CERT, 'utf8'),
      //   ca: fs.readFileSync(process.env.REACT_HTTPS_CA, 'utf8'),
      //   passphrase: process.env.REACT_HTTPS_PASS
      // };
      config.proxy = {
        "/api": { target:"https://laizn.com",secure: false }
      };
      // Return your customised Webpack Development Server config.
      return config;
    };
  }
}

