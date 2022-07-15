module.exports = (sequelize, Sequelize) => {
    const WordpressPlugin = sequelize.define("plugin", {
        id: {
            primaryKey: true,
            type: Sequelize.STRING
        },
        name: {
            type: Sequelize.STRING
        },
        version: {
            type: Sequelize.STRING
        },
        author: {
            type: Sequelize.STRING
        },
        description: {
            type: Sequelize.STRING
        }
    });
    return WordpressPlugin;
};