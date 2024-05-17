// SPDX-License-Identifier: MIT
pragma solidity ^0.8.7;

contract PriceFeedMock {
    function latestRoundData()
    external
    pure
    returns (
        uint80 roundId,
        int answer,
        uint startedAt,
        uint updatedAt,
        uint80 answeredInRound
    )
    {
        return (0, 3087, 0, 0, 0); // Mock data, 1 ETH = 3087 USDT
    }
}
