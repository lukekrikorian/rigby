const path = require("path");
const CleanWebpackPlugin = require("clean-webpack-plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");


module.exports = {
    mode: "development",
    entry: "./client/index.jsx",
    plugins: [
        new CleanWebpackPlugin(["./static/dist"]),
        new HtmlWebpackPlugin({
            template: "static/template.html",
            inject: false,
        })
    ],
    output: {
        filename: "[name].[contentHash].js",
        path: path.resolve(__dirname + "/static/dist")
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: ["cache-loader",
                {
                    loader: 'babel-loader',
                    options: {
                        babelrc: false,
                        presets: ["@babel/env", "@babel/react"]
                    }
                }]
            },
            {
                test: /\.css$/,
                use: ["style-loader", "css-loader"]
            }
        ]
    },
    resolve: {
        alias: {
            "react": "preact-compat",
            "react-dom": "preact-compat"
        },
        extensions: [".js", ".jsx"]
    }
};