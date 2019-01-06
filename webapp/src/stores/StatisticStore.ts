import Api from "../core/Api";

const WorkflowStatusLabels = {
	1: "WorkflowCreated",
	10: "WorkflowReceiveInformation",
	20: "WorkflowWaitForPickup",
	30: "WorkflowTransporting",
	50: "WorkflowDelivered",
}

class StatisticStore {
    private static instance: StatisticStore;

    public static getInstance() : StatisticStore {
        if ( StatisticStore.instance == null ) {
            StatisticStore.instance = new StatisticStore();
        }

        return StatisticStore.instance;
    }
}

export default StatisticStore;
