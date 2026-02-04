const path = require('path');

module.exports = {
    entry: {
        signup: './src/signup.js',
        login: './src/login.js',
    },
    output: {
        filename: '[name].js', // signup.js & login.js
        path: path.resolve(__dirname, 'static/dist'),
        clean: true,
        publicPath: '/static/dist/',
    },
    module: {
        rules: [
            {
                test: /\.(png|svg|jpg|jpeg|gif)$/i,
                type: 'asset/resource',
            },
        ],
    },
};
