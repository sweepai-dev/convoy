import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CardComponent } from 'src/app/components/card/card.component';
import { FormBuilder, FormGroup, ReactiveFormsModule } from '@angular/forms';
import { CreateSubscriptionService } from '../create-subscription/create-subscription.service';
import { ButtonComponent } from 'src/app/components/button/button.component';
import { GeneralService } from 'src/app/services/general/general.service';

@Component({
	selector: 'convoy-create-subscription-filter',
	standalone: true,
	imports: [CommonModule, CardComponent, ReactiveFormsModule, ButtonComponent],
	templateUrl: './create-subscription-filter.component.html',
	styleUrls: ['./create-subscription-filter.component.scss']
})
export class CreateSubscriptionFilterComponent implements OnInit {
	@Output('filterSchema') filterSchema: EventEmitter<any> = new EventEmitter();
	subscriptionFilterForm: FormGroup = this.formBuilder.group({
		request: [],
		schema: []
	});
	constructor(private formBuilder: FormBuilder, private createSubscriptionService: CreateSubscriptionService, private generalService: GeneralService) {}

	ngOnInit(): void {}

	async testFilter() {
		if (!this.convertStringToJson(this.subscriptionFilterForm.value.schema) || !this.convertStringToJson(this.subscriptionFilterForm.value.request)) return;
		this.subscriptionFilterForm.value.schema = this.convertStringToJson(this.subscriptionFilterForm.value.schema);
		this.subscriptionFilterForm.value.request = this.convertStringToJson(this.subscriptionFilterForm.value.request);
		try {
			const response = await this.createSubscriptionService.createSubsriptionFilter(this.subscriptionFilterForm.value);
			this.generalService.showNotification({ message: response.message, style: 'success' });
			console.log(response);
		} catch (error) {
			return error;
		}
	}

	convertStringToJson(str: string) {
		try {
			const jsonObject = JSON.parse(str);
			return jsonObject;
		} catch {
			this.generalService.showNotification({ message: 'Event data is not entered in correct JSON format', style: 'error' });
			return false;
		}
	}
}
