const db = require("../models");
const PluginModel = db.plugins;

// Create and Save a new Plugin
exports.create = (req, res) => {
	// Validate request
	if (!req.body.id || !req.body.name || !req.body.version || !req.body.author || !req.body.description) {
		res.status(400).send({
			message: "All details must be provided: ID, name, version, author, description"
		});
    	return;
  	}
	
	const plugin = {
		id: req.body.id,
		name: req.body.name,
    	version: req.body.version,
		author: req.body.author,
    	description: req.body.description,
  	};
  
  	// Save Wordpress plugin in the database
  	PluginModel.create(plugin)
	.then(data => {
		res.send({
			status: "success", 
			http_status_code: 201
		});
    })
    .catch(err => {
		res.status(500).send({
			message: err.message || "Some error occurred while creating the Tutorial."
		});
    });
};

// Retrieve all Plugins from the database.
exports.findAll = (req, res) => {
	PluginModel.findAll()
	.then(data => {
		res.send({
			plugins: data,
			http_status_code: 200
		});
	})
    .catch(err => {
	res.status(500).send({
		message:
			err.message || "Some error occurred while retrieving tutorials."
  		});
	});  
};

// Find a single Plugin with an id
exports.findOne = (req, res) => {
	const id = req.params.id;
  	PluginModel.findByPk(id)
    .then(data => {
		if (data) {
        	res.send(data);
      	} else {
        	res.status(404).send({
          		message: `Cannot find plugin with id=${id}.`
        	});
      	}
    })
    .catch(err => {
		res.status(500).send({
        	message: "Error retrieving plugin with id=" + id
      	});
    });  
};

// Update a Plugin by the id in the request
exports.update = (req, res) => {
	// Validate request
	if (!req.body.id || !req.body.version) {
		res.status(400).send({
			message: "All details must be provided to update plugin: ID, version"
		});
    	return;
  	}

  	PluginModel.update(req.body, {
    	where: { id: req.body.id }
  	})
    .then(num => {
		if (num == 1) {
        	res.send({
				status: "success", 
				http_status_code: 202
        	});
      	} else {
    		res.send({
    			status: `Cannot update plugin with id=${id}. No plugin found or request body is empty!`
        	});
      	}
    })
    .catch(err => {
    	res.status(500).send({
        	message: `Erroor updating plugin with id=${id}.`
      	});
    });
};