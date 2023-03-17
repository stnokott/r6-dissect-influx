import type { matches } from "../matches";

export function createRoundInfo(won: boolean, role: "ATTACK" | "DEFENSE"): matches.RoundInfo {
	const time = new Date("2023-01-01 00:00:00")
	const roundInfo: matches.RoundInfo = {
		MatchID: "123",
		Time: time.toISOString(),
		SeasonSlug: "Y8S1",
		MatchType: "RANKED",
		GameMode: "BOMB",
		MapName: "NIGHTHAVEN_LABS",
		Teams: createTeams("FooBar", role),
		Site: "0F Basement",
		Won: won,
		WinCondition: "KILLED_OPPONENTS",
		TeamIndex: role === "ATTACK" ? 0 : 1,
		PlayerName: "FooBar"
	};
	return roundInfo;
}

const attackerOps = new Array<string>(
	"AMARU",
	"ACE",
	"KALI",
	"THERMITE",
	"THATCHER"
);
const defenderOps = new Array<string>(
	"ARUNI",
	"KAPKAN",
	"CASTLE",
	"MELUSI",
	"CLASH"
);

function createTeams(observerName: string, observerTeamRole: "ATTACK" | "DEFENSE"): [matches.Team, matches.Team] {
	const numberOfTeams = 2;
	const numberOfPlayers = 5;
	const teams: [matches.Team, matches.Team] = [
		{
			Role: "ATTACK",
			Players: new Array<matches.Player>(numberOfPlayers)
		},
		{
			Role: "DEFENSE",
			Players: new Array<matches.Player>(numberOfPlayers)
		}
	]
	for (let teamIndex = 0; teamIndex < numberOfTeams; teamIndex++) {
		for (let playerIndex = 0; playerIndex < numberOfPlayers; playerIndex++) {
			teams[teamIndex].Players[playerIndex] = {
				Username: `Player ${playerIndex + teamIndex + 1}`,
				Operator: teamIndex === 0 ? attackerOps[playerIndex] : defenderOps[playerIndex]
			};
		}
	}
	teams[observerTeamRole === "ATTACK" ? 0 : 1].Players[0].Username = observerName;
	return teams;
}
