## Opensolar on Lightning and other Layer 2 Payment Networks

The idea behind openx is that it can be expanded to other blockchains and traditional infrastructure quite easily. In this document, we explore how we might create an openx instance based off the Bitcoin Lightning Network.

Openx in itself, needs a couple things in order for it to function correctly:

1. The availability of asset creation
2. The availability of a mechanism to be able to exchange assets easily (ie a secondary market for assets)

Openx leaves it to implementations to take care of the following:

1. The storage layer of entities that are internal to the system
2. The contracts that define what kind of investments that can be made

Other stuff like commitment schemes and similar are internal to the system and do not affect the functioning of openx as a whole. For opensolar and Stellar, we use the following parameters to define the opensolar platform:

1. Assets on Stellar
2. Secondary markets based on fungible asset creation
3. Data storage enabled by ipfs
4. Running a semi trusted smart contract side and committing transaction at intervals to the blockchain

The opensolar platform is, as a result semi trusted, since the investors and recipients trust the platform to not lie about project parameters but do not trust the platform enough to allow custodianship of money. The Stellar blockchain is a good fit for this use case since:

1. The protocol's lack of incentive for validators do not allow for a decentralized setting.
2. The honesty of quorums is tied to a loss of reputation and goodwill in the community.

Firms that have a stake in Stellar are encouraged to run nodes although popular research shows that most validators on Stellar rely on 2/3 nodes for consensus and their downtime can prevent the protocol from creating more blocks (as experienced in May 2019)

Bitcoin, on the other hand is more decentralized than Stellar and hence an openx platform that builds on top of Bitcoin can aim to do more trust-less things. But

1. Asset creation in bitcoin involving coloured coins or similar is expensive and non fungible
2. If such assets are non fungible, some assets will be perceived to have more value than the other and this might not prove to be good for the platform itself?
3. On chain transactions on Bitcoin are slow and expensive compared to Stellar.

The focus as a result is shifted to looking at potential Layer 2 alternatives which can look to address some of these concerns. There are two solutions we could look towards:

1. The Lightning Network
2. Sidechains (Liquid)

We will take a look at how we could potentially build opensolar on a Lightning Network based ecosystem.

#### The existence of a hub

One of the things that isn't strictly needed but needed in order to have a better user experience around routing channels is to construct a Lightning Hub. This hub would act as the control node with respect to creating channels, routing them around, etc. The hub would be connected to an instance of a core node, which the platform would be required to run.

#### Signing up as an investor

In the event of lightning however, there are a couple things that must be done after an investor puts funds into his account:

1. Create a one-way lightning channel with the hub holding a majority of his funds in the channel - The platform responds by putting in some liquidity into the channel itself in order to potentially pay back the investor after he invests in a particular platform.

2. Create a channel of low capacity with another hub, not directly related with the platform - This is to ensure that the investor does not trust the platform even for their investment and in the case the platform acts rogue, they can prove that the platform went rogue and get their funds back.

The above two conditions assume that the user has bitcoin with them and would like to invest in the project. Another possibility is:

1. The users signs up directly with layer 2 liquidity using something like Thor and route payments to the lightning Hub for investments.

All assumptions hold as above but this would mean that these investors need not sign up on openx - they could directly route to the lightning hub provided that regulations allow for such investment opportunities.

#### Investing in a project

When an investor wants to invest in a project, he can:

1. Push funds in their existing channel with the lighning hub
2. Route payments to the lightning hub

1 is marginally cheaper and more effective.

#### When a project reaches the funding goal

When a project reaches its funding target, the project splices out funds from the channel into a channel involving the project escrow. This escrow would in turn have channels with recipients, contractors, developers, etc who would be able to receive funds from the platform.

#### Proof of investment in a particular project

To track investments, we could either issue a token on behalf of the platform using specs offered by [rgb](https://github.com/rgb-org/spec) on the lightning network. We could also commit a single transaction to the blockchain whenever a project reaches its funding goals as proof of investment into the particular project.

#### The Project Escrow

The project escrow would be a mini hub which would redistribute funds out to investors and other parties involved with the project. The keys to this escrow account would be controlled by a federation of entities to minimize trust in a single entity.

#### Withdrawing funds

When an entity wants to withdraw funds form the lightning channel, it can simply close its existing channel with the mini hub and create a new one (create a tx on chain where one of the outputs is a splice in transaction). Conversion from BTC to USD is easy since there are multiple services like Hodl Hodl offering trustless exchange of btc into usd.

#### Comparisons with existing model

1. Trust - the platform becomes a bit more trustless since the investor can prove that the platform went rogue and get their funds back
2. Capital Lockup - the new model requires capital lockup to create new channels and have money in lockup to payback investors in case the platform goes rogue.
3. Faster - Lightning offers immediate settlement
4. More decentralized - Lightning, even with its hub an spoke model, is more decentralized than Stellar
5. Faster fiat off / on ramps - Courtesy to Bitcoin's market advantage, Bitcoin on an off ramps are much more readily available
6. Atomic Swaps - Atomic swaps for bitcoin-like coins to replace the DEX in Stellar.
7. Anonymity - since the transactions are off chain, this offers more anonymity to investors
