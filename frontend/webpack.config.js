const path = require('path');

module.exports = {
    entry: './src/signup.js',
    output: {
        filename: 'signup.js',
        path: path.resolve(__dirname, 'static/dist'),
        clean: true,
        publicPath: '/static/dist/'
    },
    module: {
        rules: [
            {
                test: /\.(png|svg|jpg|jpeg|gif)$/i,
                type: 'asset/resource'
            }
        ]
    }
};
