
Algo

score for auctioning a piece: 
upper cases represent real state, lower case represent arbitrary weights/coefficients
    `T` = # of pieces played this round
    `N` = # of pieces of this artist played this round
    `X` = # of pieces of other artists played this round (simple sum for simplicity, could be weighted)
    `P` = # of people who hold this artist besides self
    `C` = # of pieces of this artist in self collection
    `L` = # of pieces of other artists in self collection
    `H` = # of pieces of this artist in self hand
    `R_(0,1,2,3)` = Range of values ie. [0, 10, 20, 30] or [0, 30, 40, 50].

```go
// end the round? 
if N == 4 {

    // if 2 other artists have 4, then playing this could be worth `H * 20` because it takes it from 3rd to first

    // number of pieces i boost to 1st - number of pieces i drop 
    my_delta = C - L 
    // other players' pieces i boost to 1st - other players' pieces i drop
    other_delta = (N-C) - (T-N)

    my_delta - other_delta - MAX( playing Other artist)

    // We add (X / k) where k is some constant because there is risk to not playing this piece if another artist comes in first instead

    // we subtract N-C because N-C is the number of pieces held by other players. We are boosting these players' scores by playing this piece
}
// doesnt end the round
else {
    my_delta  = C - L
    other_delta = (N-C) - (T-N)
    competitor1 = MAX(X) // max of other artists
    competitor2 = MAX(X -- competitor1) // max of other artists, excluding competitor1
    competitor3 = MAX(X -- competitor1 -- competitor2) // max of other artists, excluding competitor1 and competitor2
    competitiveness = 10 // how competitive is this artist? 
    // with tiebreakers, we want to play the most competitive artist
    // without tiebreakers: 
    
    // comp deltas are how far ahead the competitors are. includes tiebreakers. 
    comp_delta_1 = competitor1 - N
    comp_delta_2 = competitor2 - N
    comp_delta_3 = competitor3 - N

    // if comp_delta_1 is less than 10, playing piece puts the artist in first. 
    if comp_delta_1 < PointsPerPiece {
        competitiveness += 10
    }
    // if comp_delta_2 is less than 10, playing piece puts the artist in second.
    if comp_delta_2 < PointsPerPiece {
        competitiveness += 5
    }

    // if comp_delta_3 is less than 10, playing piece puts the artist in third.
    if comp_delta_3 < PointsPerPiece {
        competitiveness += 2
    }

    // if playing this card still leaves this artist in last place, it's not worth much
    if comp_delta_3 > PointsPerPiece {
        competitiveness -= 5
    }


    expected_bid = competitiveness * weighted_average(R) // how much do i expect to get from auctioning this piece?

    // how much impact does this piece have on the rankings * net impact on player payouts
    // plus expected revenue from auctioning this piece
}
    final_score = competitiveness_delta * (my_delta - other_delta) + expected_bid



```