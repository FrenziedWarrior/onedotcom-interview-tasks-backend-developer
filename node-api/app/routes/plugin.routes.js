module.exports = app => {
    const plugins = require("../controllers/plugin.controller.js");
    var router = require("express").Router();

    // Create a new Wordpress plugin
    router.post("/plugins", plugins.create);

    // Retrieve all Wordpress plugins
    router.get("/plugins", plugins.findAll);

    // Retrieve a single Wordpress plugin with id
    router.get("/plugins/:id", plugins.findOne);

    // Update a Wordpress plugin with id
    router.patch("/plugins", plugins.update);

    app.use('/api/v1/wordpress', router);
};