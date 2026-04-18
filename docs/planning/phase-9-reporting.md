# Phase 9: Reporting & Analytics

## Goal
Transform raw data from transactions and inventory into actionable business insights.

## High-Level Requirements
- **Sales Reports**: Daily, weekly, monthly sales breakdown.
- **Product Mix (PMix)**: Analysis of menu items performance.
- **COGS & Profitability**: Analysis of net profit vs material usage cost.
- **Low Stock Alerts**: Automated reporting/notifications for items nearing `min_stock`.
- **Top Suppliers**: Analysis of purchasing spend.

## Proposed Components
- **Aggregator Service**: A dedicated service or materialized views in DB to calculate summaries.
- **Export Engine**: CSV/PDF generation for reports.

## Key Metrics to Track
- **Gross Profit**: Total Revenue - Total COGS.
- **Average Order Value (AOV)**: Revenue / Total Orders.
- **Waste Management**: Stock movements marked as "DEDUCTION/Wasted".

## Execution Steps
- [ ] Implement Sales aggregation API.
- [ ] Create COGS calculation engine.
- [ ] Build "Top Selling" and "Worst Selling" menu reports.
- [ ] Implement Excel/CSV Export for financial auditing.
- [ ] Implement Material Consumption report based on recipes.
