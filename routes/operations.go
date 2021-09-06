package routes

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/med8bra/moni-api-go/middleware"
	"github.com/med8bra/moni-api-go/models"
)

// const multer = require("multer");

// //middleware for file upload
// const upload = multer({
//   limits: {
//     fileSize: 5000000, // max file size 1MB = 1000000 bytes
//   },
//   fileFilter(req, file, cb) {
//     console.log("file...", file);
//     if (!file.originalname.match(/\.(jpeg|jpg|png)$/)) {
//       cb(new Error("only upload files with jpg or jpeg format."));
//     }
//     cb(undefined, true); // continue with upload
//   },
// });

// @route    GET api/operations/:id
// @desc     Get a Operations
// @access   Private
func getOperationsId(c *gin.Context) {
	operationId := c.Param("id")
	var operation models.Operation

	if err := models.Database.First(&operation, operationId).Error; err != nil {
		c.JSON(http.StatusNotFound, "operation not found")
		return
	}
	c.JSON(http.StatusOK, operation)
}

// @route    GET api/operations/commission/agent/id
// @desc     Get all Operations commission for agent
// @access   Private
func getOperationsCommissionAgentId(c *gin.Context) {
	agentId := c.Param("id")
	var operations []models.Operation

	if err := models.Database.
		Where("agent_id = ?", agentId).
		Order("date DESC").
		Find(&operations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	var dollar, soles float64

	for _, operation := range operations {
		if operation.Transaction.CurrencyTo == "Dollars" && operation.Transaction.Status == "Finalizado" {
			dollar += operation.Transaction.AmountToPay
		} else if operation.Transaction.CurrencyTo == "Soles" && operation.Transaction.Status == "Finalizado" {
			dollar += operation.Transaction.AmountToPay
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"dollar": dollar,
		"soles":  soles,
	})
}

// @route    GET api/operations/commission/admin
// @desc     Get all Operations commission
// @access   Private
func getOperationsCommissionAdmin(c *gin.Context) {
	var operations []models.Operation

	if err := models.Database.
		Order("date DESC").
		Find(&operations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	var dollar, soles float64

	for _, operation := range operations {
		if operation.Transaction.CurrencyTo == "Dollars" && operation.Transaction.Status == "Finalizado" {
			dollar += operation.Transaction.AmountToPay
		} else if operation.Transaction.CurrencyTo == "Soles" && operation.Transaction.Status == "Finalizado" {
			dollar += operation.Transaction.AmountToPay
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"dollar": dollar,
		"soles":  soles,
	})
}

// @route    GET api/operations/agent/id
// @desc     Get all Operations
// @access   Private
func getOperationsAgentId(c *gin.Context) {
	agentId := c.Param("id")
	var payload struct {
		Page  uint `form:"page"`
		Limit uint `form:"limit"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	offset := (payload.Page - 1) * payload.Limit

	var operations []models.Operation
	if err := models.Database.
		Where("agentId = ?", agentId).
		Offset(int(offset)).
		Limit(int(payload.Limit)).
		Find(&operations).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, operations)
}

// @route    GET api/operations
// @desc     Get all Operations
// @access   Private
func getOperations(c *gin.Context) {
	var payload struct {
		Page  uint `form:"page"`
		Limit uint `form:"limit"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	offset := (payload.Page - 1) * payload.Limit

	var operations []models.Operation
	if err := models.Database.
		Offset(int(offset)).
		Limit(int(payload.Limit)).
		Find(&operations).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, operations)
}

// @route    GET api/operations/count
// @desc     Get count of operations
// @access   Private
func getOperationsCount(c *gin.Context) {
	var operationsCount int64
	if err := models.Database.
		Model(&models.Operation{}).
		Count(&operationsCount).
		Error; err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, operationsCount)
}

// @route    GET api/operations/photo/:id
// @desc     Get image of user operation
// @access   Public
func getOperationsPhotoId(c *gin.Context) {
	operationId := c.Param("id")
	var operation models.Operation
	if err := models.Database.
		Preload("UserTransactionPhoto").
		First(&operation, operationId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	image := bytes.NewBuffer(operation.UserTransactionPhoto.Content)
	c.DataFromReader(http.StatusOK,
		int64(image.Len()),
		"image/jpeg",
		image,
		map[string]string{
			"Content-Disposition": `attachment; filename="userTransactionPhoto.jpeg"`,
		},
	)
}

// @route    GET api/operations/user
// @desc     Get all Operations of single User
// @access   Private
func getOperationsUser(c *gin.Context) {
	userId := c.MustGet("user").(*models.AuthUser).ID
	var payload struct {
		Page  uint `form:"page"`
		Limit uint `form:"limit"`
	}
	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	offset := (payload.Page - 1) * payload.Limit

	var operations []models.Operation
	if err := models.Database.
		Where(models.Account{UserID: userId}).
		Offset(int(offset)).
		Limit(int(payload.Limit)).
		Find(&operations).Error; err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, "server error")
		return
	}
	c.JSON(http.StatusOK, operations)
}

// @route    POST api/operations
// @desc     Add a new operation
// @access   Private
// func postOperations(c *gin.Context) {
// 	var payload
// }

// router.post(
//   "/",
//   [
//     auth,
//     upload.single("file"),
//     [
//       check("profileDetails", "Profile Details are required").not().isEmpty(),
//       check("bankDetails", "Bank Account Details are required").not().isEmpty(),
//       check("destinationBank", "Destination Bank is required").not().isEmpty(),
//       check("transaction", "Transaction Details Required").not().isEmpty(),
//       check("agentBank", "Agent Bank Details Required").not().isEmpty(),
//     ],
//   ],
//   async (req, res) => {
//     const errors = validationResult(req);
//     if (!errors.isEmpty()) {
//       return res.status(400).json({ errors: errors.array() });
//     }

//     let {
//       profileDetails,
//       bankDetails,
//       transaction,
//       agent,
//       commissionType,
//       commissionValue,
//       destinationBank,
//       transactionNumber,
//       agentBank,
//       savings,
//       exchange,
//     } = req.body;

//     profileDetails = JSON.parse(profileDetails);
//     bankDetails = JSON.parse(bankDetails);
//     transaction = JSON.parse(transaction);
//     agent = JSON.parse(agent);
//     commissionType = JSON.parse(commissionType);
//     commissionValue = JSON.parse(commissionValue);
//     destinationBank = JSON.parse(destinationBank);
//     transactionNumber = JSON.parse(transactionNumber);
//     savings = JSON.parse(savings);
//     exchange = JSON.parse(exchange);
//     agentBank = JSON.parse(agentBank);

//     let userTransactionPhoto = "";

//     if (req.file) {
//       userTransactionPhoto = req.file.buffer;
//     }

//     try {
//       let user = await Agent.findById(agent);
//       console.log(user);
//       const newOperation = new Operation({
//         profileDetails,
//         bankDetails,
//         transaction,
//         destinationBank,
//         transactionNumber,
//         user: req.user.id,
//         agent,
//         agentName: user.name,
//         agentEmail: user.email,
//         savings,
//         exchange,
//         agentBank,
//         userTransactionPhoto,
//       });

//       const operation = await newOperation.save();

//       console.log("operation: ---", newOperation);
//       res.json(newOperation);

//       if (!user)
//         return res.status(404).json({ msg: "Agent with that ID not found" });

//       user = await Agent.findByIdAndUpdate(
//         agent,
//         { $push: { operations: operation._id } },
//         { new: true }
//       );
//     } catch (err) {
//       console.error(err.message);
//       res.status(500).send("Server Error");
//     }
//   }
// );

// // @route    PUT api/operations/:id
// // @desc     Update an operation
// // @access   Private
// router.put("/:id", auth, async (req, res) => {
//   const errors = validationResult(req);
//   if (!errors.isEmpty())
//     return res.status(400).json({ errors: errors.array() });

//   const {
//     profileDetails,
//     bankDetails,
//     destinationBank,
//     transactionNumber,
//     transaction,
//     agentBank,
//     userTransactionPhoto,
//     agentTransactionNumber,
//   } = req.body;

//   // Build operation object
//   const accountFields = {};
//   if (profileDetails) accountFields.profileDetails = profileDetails;
//   if (bankDetails) accountFields.bankDetails = bankDetails;
//   if (transaction) accountFields.transaction = transaction;
//   if (transactionNumber) accountFields.transactionNumber = transactionNumber;
//   if (agentBank) accountFields.agentBank = agentBank;
//   if (destinationBank) accountFields.destinationBank = destinationBank;
//   if (userTransactionPhoto) accountFields.userTransactionPhoto = userTransactionPhoto;
//   if (agentTransactionNumber)
//     accountFields.agentTransactionNumber = agentTransactionNumber;

//   try {
//     let operation = await Operation.findById(req.params.id);

//     if (!operation) return res.status(404).json({ msg: "Operation not found" });
//     let user = await Agent.findById(operation.agent);

//     operation = await Operation.findByIdAndUpdate(
//       req.params.id,
//       { $set: accountFields },
//       { new: true }
//     );
//     let commissionType = operation.transaction.currencyFrom,
//       commissionValue = operation.transaction.amountReceive * 0.0006;
//     if (operation.transaction.status === "Finalizado") {
//       if (!user)
//         return res.status(404).json({ msg: "Agent with that ID not found" });
//       if (commissionType === "Soles") {
//         user = await Agent.findByIdAndUpdate(
//           operation.agent,
//           { $inc: { commissionSoles: commissionValue } },
//           { new: true }
//         );
//         console.log(user);
//       } else if (commissionType === "Dollars") {
//         user = await Agent.findByIdAndUpdate(
//           operation.agent,
//           { $inc: { commissionDollars: commissionValue } },
//           { new: true }
//         );
//         console.log(user);
//       }
//     }
//     res.json(operation);
//   } catch (err) {
//     console.error(err.message);
//     res.status(500).send("Server error");
//   }
// });

// // @route    DELETE api/operations/:id
// // @desc     Delete a operation
// // @access   Private
// router.delete("/:id", auth, async (req, res) => {
//   try {
//     const operation = await Operation.findById(req.params.id);

//     if (!operation) return res.status(404).json({ msg: "Operation not found" });

//     user = await Agent.findByIdAndUpdate(
//       operation.agent,
//       { $pull: { operations: operation._id } },
//       { new: true }
//     );

//     await Operation.findByIdAndRemove(req.params.id);

//     res.json({ msg: "Operation removed" });
//   } catch (err) {
//     console.error(err.message);
//     res.status(500).send("Server error");
//   }
// });

// module.exports = router;

func operations(r *gin.RouterGroup) {
	// private
	r.Use(middleware.Auth).
		GET("/commission/agent/:id", getOperationsCommissionAgentId).
		GET("/commission/admin", getOperationsCommissionAdmin).
		GET("/agent/:id", getOperationsAgentId).
		GET("/photo/:id", getOperationsPhotoId).
		GET("/count", getOperationsCount).
		GET("/user", getOperationsUser).
		GET("/:id", getOperationsId).
		GET("", getOperations)
}
