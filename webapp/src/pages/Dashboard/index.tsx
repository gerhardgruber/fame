import * as React from 'react';
import { observer } from 'mobx-react';
import { Icon, Row, Col, Result, Spin } from 'antd'
import Page from '../../components/Page';
import UiStore from '../../stores/UiStore';
import Api from '../../core/Api';
import { observable } from 'mobx';

const uiStore = UiStore.getInstance( );

@observer
export default class Dashboard extends Page {
  @observable dates: number;

  @observable positive: number;

  @observable present: number;

  private static _green = "green";
  private static _orange = "#FF9100";
  private static _red = "red";
  private static _yellow = "#FFEA00";

  componentDidMount() {
    Api.GET("/dashboard/status", {}).then((response) => {
      this.dates = response.data.data.Dates;
      this.positive = response.data.data.Positive;
      this.present = response.data.data.Present;
    })
  }

  pageTitle(): string {
    return "DASHBOARD";
  }

  renderResult(type: string) {
    let value: number;
    if ( type === "present" ) {
      value = this.present;
    } else {
      value = this.positive;
    }

    let factor = 0;
    if (this.dates !== 0 ) {
      factor = value / this.dates;
    }

    let icon: string;
    let color: string;
    if ( factor >= 0.75 || this.dates === 0 ) {
      icon = "check-circle";
      color = Dashboard._green;
    } else if ( factor >= 0.5 ) {
      icon = "warning";
      color = Dashboard._yellow;
    } else if ( factor >= 0.25 ) {
      icon = "warning";
      color = Dashboard._orange;
    } else {
      icon = "close-circle";
      color = Dashboard._red;
    }

    return <Result
        icon={<Icon type={icon} twoToneColor={color} theme="twoTone" />}
        title={`${value} / ${this.dates}`}
      />
  }

  renderContent() {
    if (this.dates === undefined) {
      return <Spin />;
    }

    return  <div>
              <Row>
                <Col md={12}>
                  <h2 style={{textAlign: 'center'}}>{uiStore.T('DASHBOARD_FEEDBACK')}</h2>
                </Col>
                <Col md={12}>
                  <h2 style={{textAlign: 'center'}}>{uiStore.T('DASHBOARD_PRESENT')}</h2>
                </Col>
              </Row>
              <Row>
                <Col md={12}>
                  {this.renderResult("positive")}
                </Col>
                <Col md={12}>
                  {this.renderResult("present")}
                </Col>
              </Row>
              <Row>
                <Col md={24} style={{padding: '50px', paddingBottom: '10px'}}>
                  {uiStore.T( "DASHBOARD_LEGEND" )}:
                </Col>
              </Row>
              <Row>
                <Col md={24} style={{paddingLeft: '50px'}}>
                  <Icon type="check-circle" theme="twoTone" twoToneColor={Dashboard._green} style={{fontSize: "2rem", marginRight: '0.5rem'}} />
                  <Icon type="warning" theme="twoTone" twoToneColor={Dashboard._yellow} style={{fontSize: "2rem", marginRight: '0.5rem'}} />
                  <Icon type="warning" theme="twoTone" twoToneColor={Dashboard._orange} style={{fontSize: "2rem", marginRight: '0.5rem'}} />
                  <Icon type="close-circle" theme="twoTone" twoToneColor={Dashboard._red} style={{fontSize: "2rem", marginRight: '0.5rem'}} />
                </Col>
              </Row>
            </div>
  }
}