package sandbox

// sandbox contains sandbox data that can be used to spin up openx for presentations at demos in conjunction
// with the openx-frontend repo. Can be edited to one's needs. File last updated: May 2019
import (
	"encoding/json"
	"io/ioutil"
	"log"

	xlm "github.com/YaleOpenLab/openx/chains/xlm"
	assets "github.com/YaleOpenLab/openx/chains/xlm/assets"
	wallet "github.com/YaleOpenLab/openx/chains/xlm/wallet"
	consts "github.com/YaleOpenLab/openx/consts"
	// database "github.com/YaleOpenLab/openx/database"
	opensolar "github.com/YaleOpenLab/openx/platforms/opensolar"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// parseYamlProject reparses yaml for an existing project
func parseYamlProject(fileName string, feJson string, projIndex int) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./data")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	project, err := opensolar.RetrieveProject(projIndex)
	if err != nil {
		return err
	}

	termsHelper := viper.Get("Terms").(map[string]interface{})
	if projIndex == 8 {
		terms := make([]opensolar.TermsHelper, 9)
		i := 0
		for _, elem := range termsHelper {
			// elem inside here is a map of "variable": values.
			newMap := elem.(map[string]interface{})
			terms[i].Variable = newMap["variable"].(string)
			terms[i].Value = newMap["value"].(string)
			terms[i].RelevantParty = newMap["relevantparty"].(string)
			terms[i].Note = newMap["note"].(string)
			terms[i].Status = newMap["status"].(string)
			terms[i].SupportDoc = newMap["supportdoc"].(string)
			i += 1
		}

		project.Terms = terms
	} else {
		terms := make([]opensolar.TermsHelper, 6)
		i := 0
		for _, elem := range termsHelper {
			// elem inside here is a map of "variable": values.
			newMap := elem.(map[string]interface{})
			terms[i].Variable = newMap["variable"].(string)
			terms[i].Value = newMap["value"].(string)
			terms[i].RelevantParty = newMap["relevantparty"].(string)
			terms[i].Note = newMap["note"].(string)
			terms[i].Status = newMap["status"].(string)
			terms[i].SupportDoc = newMap["supportdoc"].(string)
			i += 1
		}

		project.Terms = terms
	}

	var executiveSummary opensolar.ExecutiveSummaryHelper

	execSummaryReader := viper.Get("ExecutiveSummary.Investment").(map[string]interface{})
	execSummaryWriter := make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Investment = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.Financials").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Financials = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.ProjectSize").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.ProjectSize = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.SustainabilityMetrics").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.SustainabilityMetrics = execSummaryWriter

	project.ExecutiveSummary = executiveSummary

	var bullets opensolar.BulletHelper
	bullets.Bullet1 = viper.Get("Bullets.Bullet1").(string)
	bullets.Bullet2 = viper.Get("Bullets.Bullet2").(string)
	bullets.Bullet3 = viper.Get("Bullets.Bullet3").(string)

	project.Bullets = bullets

	var architecture opensolar.ArchitectureHelper

	architecture.SolarArray = viper.Get("Architecture.SolarArray").(string)
	architecture.DailyAvgGeneration = viper.Get("Architecture.DailyAvgGeneration").(string)
	architecture.System = viper.Get("Architecture.System").(string)
	architecture.InverterSize = viper.Get("Architecture.InverterSize").(string)

	project.Architecture = architecture

	project.Index = viper.Get("Index").(int)
	project.Name = viper.Get("Name").(string)
	project.State = viper.Get("State").(string)
	project.Country = viper.Get("Country").(string)
	project.TotalValue = viper.Get("TotalValue").(float64)
	project.Metadata = viper.Get("Metadata").(string)
	project.PanelSize = viper.Get("PanelSize").(string)
	project.PanelTechnicalDescription = viper.Get("PanelTechnicalDescription").(string)
	project.Inverter = viper.Get("Inverter").(string)
	project.ChargeRegulator = viper.Get("ChargeRegulator").(string)
	project.ControlPanel = viper.Get("ControlPanel").(string)
	project.CommBox = viper.Get("CommBox").(string)
	project.ACTransfer = viper.Get("ACTransfer").(string)
	project.SolarCombiner = viper.Get("SolarCombiner").(string)
	project.Batteries = viper.Get("Batteries").(string)
	project.IoTHub = viper.Get("IoTHub").(string)
	project.Rating = viper.Get("Rating").(string)
	project.EstimatedAcquisition = viper.Get("EstimatedAcquisition").(int)
	project.BalLeft = viper.Get("BalLeft").(float64)
	project.InterestRate = viper.Get("InterestRate").(float64)
	project.Tax = viper.Get("Tax").(string)
	project.DateInitiated = viper.Get("DateInitiated").(string)
	project.DateFunded = viper.Get("DateFunded").(string)
	project.AuctionType = viper.Get("AuctionType").(string)
	project.InvestmentType = viper.Get("InvestmentType").(string)
	project.PaybackPeriod = viper.Get("PaybackPeriod").(int)
	project.Stage = viper.Get("Stage").(int)
	project.SeedInvestmentFactor = viper.Get("SeedInvestmentFactor").(float64)
	project.SeedInvestmentCap = viper.Get("SeedInvestmentCap").(float64)
	project.ProposedInvestmentCap = viper.Get("ProposedInvestmentCap").(float64)
	project.SelfFund = viper.Get("SelfFund").(float64)
	project.SecurityIssuer = viper.Get("SecurityIssuer").(string)
	project.BrokerDealer = viper.Get("BrokerDealer").(string)
	project.EngineeringLayoutType = viper.Get("EngineeringLayoutType").(string)
	project.MapLink = viper.Get("MapLink").(string)

	project.FEText, err = parseJsonText(feJson)
	if err != nil {
		log.Fatal(err)
	}

	return project.Save()
}

// parseYaml parses yaml, creates a new project and saves it to the database
func parseYaml(fileName string, feJson string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath("./data")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading values from config file")
	}

	var project opensolar.Project
	terms := make([]opensolar.TermsHelper, 6)
	termsHelper := viper.Get("Terms").(map[string]interface{})

	i := 0
	for _, elem := range termsHelper {
		// elem inside here is a map of "variable": values.
		newMap := elem.(map[string]interface{})
		terms[i].Variable = newMap["variable"].(string)
		terms[i].Value = newMap["value"].(string)
		terms[i].RelevantParty = newMap["relevantparty"].(string)
		terms[i].Note = newMap["note"].(string)
		terms[i].Status = newMap["status"].(string)
		terms[i].SupportDoc = newMap["supportdoc"].(string)
		i += 1
	}

	project.Terms = terms
	var executiveSummary opensolar.ExecutiveSummaryHelper

	execSummaryReader := viper.Get("ExecutiveSummary.Investment").(map[string]interface{})
	execSummaryWriter := make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Investment = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.Financials").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.Financials = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.ProjectSize").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.ProjectSize = execSummaryWriter

	execSummaryReader = viper.Get("ExecutiveSummary.SustainabilityMetrics").(map[string]interface{})
	execSummaryWriter = make(map[string]string)
	for key, elem := range execSummaryReader {
		execSummaryWriter[key] = elem.(string)
	}
	executiveSummary.SustainabilityMetrics = execSummaryWriter

	project.ExecutiveSummary = executiveSummary

	var bullets opensolar.BulletHelper
	bullets.Bullet1 = viper.Get("Bullets.Bullet1").(string)
	bullets.Bullet2 = viper.Get("Bullets.Bullet2").(string)
	bullets.Bullet3 = viper.Get("Bullets.Bullet3").(string)

	project.Bullets = bullets

	var architecture opensolar.ArchitectureHelper

	architecture.SolarArray = viper.Get("Architecture.SolarArray").(string)
	architecture.DailyAvgGeneration = viper.Get("Architecture.DailyAvgGeneration").(string)
	architecture.System = viper.Get("Architecture.System").(string)
	architecture.InverterSize = viper.Get("Architecture.InverterSize").(string)

	project.Architecture = architecture

	project.Index = viper.Get("Index").(int)
	project.Name = viper.Get("Name").(string)
	project.State = viper.Get("State").(string)
	project.Country = viper.Get("Country").(string)
	project.TotalValue = viper.Get("TotalValue").(float64)
	project.Metadata = viper.Get("Metadata").(string)
	project.PanelSize = viper.Get("PanelSize").(string)
	project.PanelTechnicalDescription = viper.Get("PanelTechnicalDescription").(string)
	project.Inverter = viper.Get("Inverter").(string)
	project.ChargeRegulator = viper.Get("ChargeRegulator").(string)
	project.ControlPanel = viper.Get("ControlPanel").(string)
	project.CommBox = viper.Get("CommBox").(string)
	project.ACTransfer = viper.Get("ACTransfer").(string)
	project.SolarCombiner = viper.Get("SolarCombiner").(string)
	project.Batteries = viper.Get("Batteries").(string)
	project.IoTHub = viper.Get("IoTHub").(string)
	project.Rating = viper.Get("Rating").(string)
	project.EstimatedAcquisition = viper.Get("EstimatedAcquisition").(int)
	project.BalLeft = viper.Get("BalLeft").(float64)
	project.InterestRate = viper.Get("InterestRate").(float64)
	project.Tax = viper.Get("Tax").(string)
	project.DateInitiated = viper.Get("DateInitiated").(string)
	project.DateFunded = viper.Get("DateFunded").(string)
	project.AuctionType = viper.Get("AuctionType").(string)
	project.InvestmentType = viper.Get("InvestmentType").(string)
	project.PaybackPeriod = viper.Get("PaybackPeriod").(int)
	project.Stage = viper.Get("Stage").(int)
	project.SeedInvestmentFactor = viper.Get("SeedInvestmentFactor").(float64)
	project.SeedInvestmentCap = viper.Get("SeedInvestmentCap").(float64)
	project.ProposedInvestmentCap = viper.Get("ProposedInvestmentCap").(float64)
	project.SelfFund = viper.Get("SelfFund").(float64)
	project.SecurityIssuer = viper.Get("SecurityIssuer").(string)
	project.BrokerDealer = viper.Get("BrokerDealer").(string)
	project.EngineeringLayoutType = viper.Get("EngineeringLayoutType").(string)
	project.MapLink = viper.Get("MapLink").(string)

	project.FEText, err = parseJsonText(feJson)
	if err != nil {
		log.Fatal(err)
	}

	return project.Save()
}

// populateStaticData populates static data for the demo projects
func populateStaticData() error {
	var err error
	log.Println("populating db with static data")
	err = createAllStaticEntities()
	if err != nil {
		return err
	}
	err = parseYaml("1kwy", "data/1kw.json")
	if err != nil {
		return err
	}
	// project: One Kilowatt Project / STAGE 7 - Puerto Rico
	err = populateStaticData1kw()
	if err != nil {
		return err
	}
	err = parseYaml("1mwy", "data/1mw.json")
	if err != nil {
		return err
	}
	// project: One Megawatt Project / STAGE 4 - New Hampshire
	err = populateStaticData1mw()
	if err != nil {
		return err
	}
	err = parseYaml("10kwy", "data/10kw.json")
	if err != nil {
		return err
	}
	// project: Ten Kilowatt Project / STAGE 8 - Connecticut Homeless Shelter
	err = populateStaticData10kw()
	if err != nil {
		return err
	}
	err = parseYaml("10mwy", "data/10mw.json")
	if err != nil {
		return err
	}
	// project: Ten Megawatt Project / STAGE 2 - Puerto Rico Public School Bond
	err = populateStaticData10MW()
	if err != nil {
		return err
	}
	err = parseYaml("100kwy", "data/100kw.json")
	if err != nil {
		return err
	}
	// project: One Hundred Kilowatt Project / STAGE 1 - Rwanda Project
	err = populateStaticData100KW()
	if err != nil {
		return err
	}
	return nil
}

func populateDynamicData() error {
	var err error
	// we ignore errors here since they are bound to happen (guarantor related errors)
	err = populateDynamicData1kw()
	if err != nil {
		log.Println("error while populating 1kw project", err)
	}
	err = populateDynamicData1mw()
	if err != nil {
		log.Println("error while populating 1mw project", err)
	}
	err = populateDynamicData10kw()
	if err != nil {
		log.Println("error while populating 10kw project", err)
	}
	return nil
}

// CreateSandbox populates the database with test values
func CreateSandbox() error {
	var err error
	err = populateStaticData()
	if err != nil {
		return err
	}
	err = populateDynamicData()
	if err != nil {
		return err
	}
	err = populateAdditionalData()
	if err != nil {
		return err
	}
	return nil
}

// parseJsonText is a helper function that reads json for the frontend data
func parseJsonText(fileName string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}

// seed additional data for a few specific investors that are useful for showing in demos
func populateAdditionalData() error {
	openlab, err := opensolar.RetrieveInvestor(46)
	if err != nil {
		return err
	}
	openlab.U.Email = "martin.wainstein@yale.edu"
	openlab.U.Address = "254 Elm Street"
	openlab.U.Country = "US"
	openlab.U.City = "New Haven"
	openlab.U.ZipCode = "06511"
	openlab.U.RecoveryPhone = "1800SECRETS"
	openlab.U.Description = "The Yale OPen Lab is the open innovation lab at the Tsai Centre for Innovative Thinking at Yale"
	err = openlab.U.Save()
	if err != nil {
		return err
	}
	err = openlab.Save()
	if err != nil {
		return err
	}

	// insert data for one specific recipient
	pasto, err := opensolar.RetrieveRecipient(47)
	if err != nil {
		return err
	}
	pasto.U.Email = "supasto2018@gmail.com"
	pasto.U.Address = "Puerto Rico, PR"
	pasto.U.Country = "US"
	pasto.U.City = "Puerto Rico"
	pasto.U.ZipCode = "00909"
	pasto.U.RecoveryPhone = "1800SECRETS"
	pasto.U.Description = "S.U. Pasto School is a school in Puerto Rico"
	err = pasto.U.Save()
	if err != nil {
		return err
	}
	err = pasto.Save()
	if err != nil {
		return err
	}
	dci, err := opensolar.RetrieveEntity(1)
	if err != nil {
		return err
	}

	dci.U.Email = "dci@mit.edu"
	dci.U.Address = "MIT Media Lab"
	dci.U.Country = "US"
	dci.U.City = "Cambridge"
	dci.U.ZipCode = "02142"
	dci.U.RecoveryPhone = "1800SECRETS"
	dci.U.Description = "The Digital Currency Initiative at the MIT Media Lab"

	err = dci.U.Save()
	if err != nil {
		return err
	}

	// we now need to register the dci as an investor as well
	var inv opensolar.Investor
	inv.U = dci.U
	err = inv.Save()
	if err != nil {
		return err
	}
	var recp opensolar.Recipient
	recp.U = dci.U
	err = recp.Save()
	if err != nil {
		return err
	}

	err = xlm.GetXLM(dci.U.StellarWallet.PublicKey)
	if err != nil {
		return err
	}

	seed, err := wallet.DecryptSeed(recp.U.StellarWallet.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend
	txhash, err := assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for dci trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.StablecoinCode, recp.U.StellarWallet.PublicKey, 600, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for dci getting stableUSD: ", txhash)

	recp, err = opensolar.RetrieveRecipient(47)
	if err != nil {
		return err
	}

	seed, err = wallet.DecryptSeed(recp.U.StellarWallet.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend
	txhash, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for pasto school trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.StablecoinCode, recp.U.StellarWallet.PublicKey, 600, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for pasto school getting stableUSD: ", txhash)

	err = xlm.GetXLM(recp.U.SecondaryWallet.PublicKey)
	if err != nil {
		return err
	}

	seed, err = wallet.DecryptSeed(recp.U.SecondaryWallet.EncryptedSeed, "x")
	if err != nil {
		return err
	}

	// send the pasto school account some money so we can demo using it on the frontend
	txhash, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, seed)
	if err != nil {
		return err
	}
	log.Println("TX HASH for pasto school sec wallet trusting stableUSD: ", txhash)

	_, txhash, err = assets.SendAssetFromIssuer(consts.StablecoinCode, recp.U.SecondaryWallet.PublicKey, 10000, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Println("SEED: ", consts.StablecoinSeed)
		return err
	}
	log.Println("TX HASH for pasto school sec wallet getting stableUSD: ", txhash)

	investor1, invSeed, err := bootstrapInvestor("medici@test.com", "Medici Ventures")
	if err != nil {
		log.Fatal(err)
	}

	_, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, invSeed)
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = assets.SendAssetFromIssuer(consts.StablecoinCode, investor1.U.StellarWallet.PublicKey, 1000000, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	investor1.U.Email = "medici@test.com"
	investor1.U.Address = "101 Medici Rd"
	investor1.U.Country = "US"
	investor1.U.City = "Salt Lake City"
	investor1.U.ZipCode = "08404"
	investor1.U.RecoveryPhone = "1800SECRETS"

	err = investor1.Save()
	if err != nil {
		log.Fatal(err)
	}

	recp1, recpSeed, err := bootstrapRecipient("rwandaenergy@test.com", "Rwanda Village Energy Collective")
	if err != nil {
		log.Fatal(err)
	}

	_, err = assets.TrustAsset(consts.StablecoinCode, consts.StablecoinPublicKey, 10000000000, recpSeed)
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = assets.SendAssetFromIssuer(consts.StablecoinCode, recp1.U.StellarWallet.PublicKey, 1000000, consts.StablecoinSeed, consts.StablecoinPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	recp1.U.Email = "rwandaenergy@test.com"
	recp1.U.Address = "45 Cyangugu Rd, Rwanda"
	recp1.U.Country = "Rwanda"
	recp1.U.City = "Rusizi"
	recp1.U.ZipCode = "6502"
	recp1.U.RecoveryPhone = "1800SECRETS"

	err = recp1.Save()
	if err != nil {
		log.Fatal(err)
	}

	project, err := opensolar.RetrieveProject(8)
	if err != nil {
		log.Fatal(err)
	}

	project.RecipientIndex = 64

	err = project.Save()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
