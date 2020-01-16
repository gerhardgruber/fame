import * as React from 'react';
import { observer } from 'mobx-react';
import { Table, Button, Switch, Icon } from 'antd'
import Page from '../../components/Page';
import UiStore from '../../stores/UiStore';
import DateStore from '../../stores/DateStore';
import {Link, Redirect} from 'react-router-dom';
import { DateModel } from '../../stores/DateStore/Date';
import moment from 'moment';
import { RightType } from '../../stores/User';

const dateStore = DateStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Dates extends Page {
  state = {
    navigateTo: null,
    showPastDates: false
  };

  columns: Table<DateModel>['props']['columns'] = [ {
    title: uiStore.T( 'DATES_STATUS' ),
    dataIndex: 'ID',
    width: '8%',
    render: (id: number, dt: DateModel) => {
      const myFeedback = dt.getMyFeedback();
      if (myFeedback && myFeedback.Feedback === uiStore.dateFeedbackTypes["Yes"]) {
        return <Icon style={{color: 'green'}} type="check-circle" />;
      } else if (myFeedback && myFeedback.Feedback === uiStore.dateFeedbackTypes["No"]) {
        return <Icon style={{color: 'red'}} type="close-circle" />;
      } else {
        return <Icon style={{color: 'orange'}} type="warning" />;
      }
    }
  }, {
    title: uiStore.T( 'DATES_TITLE' ),
    dataIndex: 'Title',
  }, {
    title: uiStore.T( 'DATES_START_TIME' ),
    dataIndex: 'StartTime',
    render: (dt) => {
      return moment(dt).format(uiStore.T('DATE_TIME_FORMAT'))
    }
  }, {
    title: uiStore.T( 'DATES_END_TIME' ),
    dataIndex: 'EndTime',
    render: (dt) => {
      return moment(dt).format(uiStore.T('DATE_TIME_FORMAT'))
    }
  } ];

  componentDidMount( ) {
    this.loadDates();
  }

  loadDates() {
    dateStore.loadDates(
      this.state.showPastDates
    );
  }

  pageTitle(): string {
    return "DATES";
  }

  rowClicked = (record: DateModel) => {
    this.setState({
      navigateTo: record.ID
    });
  }

  renderButtons(): JSX.Element {
    if (!uiStore.isAdmin()) {
      return null;
    }

    return <Link to="/dates/new"><Button>
      {uiStore.T( 'DATES_ADD_DATE' )}
    </Button></Link>
  }

  renderContent(): JSX.Element {
    return <div>
             <div style={{marginBottom: '1rem'}}>
               <Switch
                 checked={this.state.showPastDates}
                 onChange={(value) => {
                   this.setState({showPastDates: value}, () => this.loadDates())
                 }}
                 />
               <span style={{marginLeft: '0.5rem'}}>
                 {uiStore.T('SHOW_PAST_DATES')}
               </span>
             </div>
             <Table
               columns={this.columns}
               dataSource={dateStore.dates}
               size={"small"}
               pagination={false}
               onRowClick={this.rowClicked} />
              {this.state.navigateTo ? <Redirect push to={"/dates/" + this.state.navigateTo} /> : null}
           </div>;
  }
}