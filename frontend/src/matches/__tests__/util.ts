import type { Player, RoundInfo, Team } from "../../game";

export function createRoundInfo(won: boolean, role: "Attack" | "Defense"): RoundInfo {
	const time = new Date("2023-01-01 00:00:00")
	const roundInfo: RoundInfo = {
		MatchID: "123",
		RoundIndex: 0,
		Time: time.toISOString(),
		SeasonSlug: "Y8S1",
		MatchType: "Ranked",
		GameMode: "Bomb",
		MapName: "NighthavenLabs",
		Teams: createTeams("FooBar", role),
		Site: "0F Basement",
		Won: won,
		WinCondition: "KilledOpponents",
		TeamIndex: role === "Attack" ? 0 : 1,
		PlayerName: "FooBar"
	};
	return roundInfo;
}

const attackerOps = new Array<string>(
	"Amaru",
	"Ace",
	"Kali",
	"Thermite",
	"Thatcher"
);
const defenderOps = new Array<string>(
	"Aruni",
	"Kapkan",
	"Castle",
	"Melusi",
	"Clash"
);

function createTeams(observerName: string, observerTeamRole: "Attack" | "Defense"): [Team, Team] {
	const numberOfTeams = 2;
	const numberOfPlayers = 5;
	const teams: [Team, Team] = [
		{
			Role: "Attack",
			Players: new Array<Player>(numberOfPlayers)
		},
		{
			Role: "Defense",
			Players: new Array<Player>(numberOfPlayers)
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
	teams[observerTeamRole === "Attack" ? 0 : 1].Players[0].Username = observerName;
	return teams;
}
