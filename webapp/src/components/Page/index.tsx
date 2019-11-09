import * as React from 'react';
import { Layout } from 'antd';
import PageHeader from '../PageHeader';

abstract class Page<P = any, S = any> extends React.Component<P, S> {
  abstract pageTitle(): string;
  abstract renderContent(): JSX.Element;

  renderButtons(): JSX.Element {
    return null;
  }

  render(): JSX.Element {
    return (
      <Layout>
          <PageHeader name={this.pageTitle()} renderButtons={this.renderButtons.bind( this )} />
          <Layout style={{padding: '10px', backgroundColor: 'white'}}>
            {this.renderContent()}
          </Layout>
      </Layout>
    )
  }
}

export default Page;