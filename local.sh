alias dollard=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .dollar
        shift
        ;;
    esac
done

if ! [ -f .dollar/data/priv_validator_state.json ]; then
  dollard init validator --chain-id "dollar-1" --home .dollar &> /dev/null

  dollard keys add validator --home .dollar --keyring-backend test &> /dev/null
  dollard genesis add-genesis-account validator 1000000ustake --home .dollar --keyring-backend test
  dollard keys add owner --recover --home .dollar --keyring-backend test --output json <<< "enjoy pen flee moral inform welcome cannon caught letter symbol patch discover bid juice toward abuse bonus gospel frame chapter magnet depart throw crater" &> /dev/null
  dollard genesis add-genesis-account owner 1000000uusdc --home .dollar --keyring-backend test
  dollard genesis add-genesis-account noble1cyyzpxplxdzkeea7kwsydadg87357qnah9s9cv 1000000uusdc --home .dollar --keyring-backend test

  dollard keys add tom --recover --home .dollar --keyring-backend test --output json <<< "dice hill prepare foam tiny album cart steel pact say never hen" &> /dev/null
  dollard genesis add-genesis-account noble1zlxkchy77rp2tmknx5n8kckntyj3wp6h6c2edm 1000000uusdn --home .dollar --keyring-backend test




  TEMP=.dollar/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.dollar.portal.owner = "noble1s7evsmath5f3ef7vk97ru2tez9k5rs00klunzu"' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.dollar.vaults.owner = "noble1s7evsmath5f3ef7vk97ru2tez9k5rs00klunzu"' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.config.chain_id = 4009' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.config.gov_chain = 1' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.config.gov_address = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ="' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json
  touch $TEMP && jq '.app_state.wormhole.guardian_sets = {"0":{"addresses":["vvpCnVfNGLf4pNkaLamrSvBdD74="],"expiration_time":0}}' .dollar/config/genesis.json > $TEMP && mv $TEMP .dollar/config/genesis.json

  dollard genesis gentx validator 1000000ustake --chain-id "dollar-1" --home .dollar --keyring-backend test &> /dev/null
  dollard genesis collect-gentxs --home .dollar &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .dollar/config/config.toml
fi

dollard start --home .dollar
