import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, Input, OnInit } from '@angular/core';

@Component({
	selector: 'convoy-badge, [convoy-badge]',
	standalone: true,
	imports: [CommonModule],
	templateUrl: './badge.component.html',
	styleUrls: ['./badge.component.scss'],
	changeDetection: ChangeDetectionStrategy.OnPush,
	host: { class: 'flex items-center' }
})
export class BadgeComponent implements OnInit {
	@Input('texture') texture: 'dark' | 'light' = 'light';
	@Input('text') text!: string;

	constructor() {}

	ngOnInit(): void {}

	get firstletters(): string {
		const firstLetters = this.text
			.split(' ')
			.map(word => word[0])
			.join('');
		return firstLetters;
	}
}
