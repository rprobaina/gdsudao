package com.example.gdsudao.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.gdsudao.R
import com.example.gdsudao.model.Area
import kotlinx.android.synthetic.main.activity_detalhes.*
import kotlin.math.roundToInt

class AreaViewHolder(inflater: LayoutInflater, parent: ViewGroup): RecyclerView.ViewHolder(inflater.inflate(R.layout.adapter_area, parent, false)){
    private var tvNome: TextView? = null
    private var tvDataPlantio: TextView? = null
    private var tvProxCorte: TextView? = null
    private var tvNumeroCortes: TextView? = null
    private var tvProgresso: TextView? = null
    private var tvDiarios: TextView? = null
    private var tvPrevisoes: TextView? = null
    private var tvNormais: TextView? = null


    init {
        tvNome = itemView.findViewById(R.id.tvNomeArea)
        tvDataPlantio = itemView.findViewById(R.id.tvDataPlantio)
        tvProxCorte = itemView.findViewById(R.id.tvDataProxCorte)
        tvNumeroCortes = itemView.findViewById(R.id.tvNumeroCortes)
        tvProgresso = itemView.findViewById(R.id.tvProgresso)
        tvDiarios = itemView.findViewById(R.id.tvDiarios)
        tvPrevisoes = itemView.findViewById(R.id.tvPrevisoes)
        tvNormais = itemView.findViewById(R.id.tvNormais)
    }

    fun bind(area: Area){
        val ST_PRIRO_CORTE = 358.0f
        val ST_OUTROS_CORTES = 281.0f

        tvNome?.text = area.nome
        tvProxCorte?.text = area.proxcorte
        tvNumeroCortes?.text = "" + area.numeroCorte

        if (!area.dataCorte.isNullOrEmpty()){
            tvDataPlantio?.text = "${area.dataCorte.substring(8, 10)}/${area.dataCorte.substring(5, 7)}/${area.dataCorte.substring(0, 4)}"
        }else{
            tvDataPlantio?.text = "ERRO"
        }


        if (!area.proxcorte.isNullOrEmpty()){
            tvProxCorte?.text = "${area.proxcorte.substring(8, 10)}/${area.proxcorte.substring(5, 7)}/${area.proxcorte.substring(0, 4)}"
            //tvProxCorte?.text = area.proxcorte.substring(0, 10).replace("-","/", false)
        }else{
            tvProxCorte?.text = "ERRO"
        }



        /*
        var dia = area.diario.toFloat().roundToInt()
        if (dia != null){
            tvDiarios?.text = dia.toString() + "%"
        }

        var pre = area.previsao.toFloat().roundToInt()
        if (pre != null){
            tvPrevisoes?.text = pre.toString() + "%"
        }

        var nor = area.normal.toFloat().roundToInt()
        if (nor != null){
            tvNormais?.text = nor.toString() + "%"
        }
        */


        if (area.diario.length > 5){
            tvDiarios?.text = area.diario.substring(0, 4) + "%" //.substring(5) + "%"
        }else{
            tvDiarios?.text = area.diario + "%"
        }

        if (area.previsao.length > 5){
            tvPrevisoes?.text = area.previsao.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvPrevisoes?.text = area.previsao + "%"  //.substring(5) + "%"
        }

        if (area.normal.length > 5){
            tvNormais?.text = area.normal.substring(0, 4) + "%"  //.substring(5) + "%"
        }else{
            tvNormais?.text = area.normal + "%"  //.substring(5) + "%"
        }



        // TODO: calcular o progresso com base nas constatnes de st e no numero de cortes
        // tvProgresso?.text = "Data Plantio: " + "|d: " + area.diario + "|p: " + area.previsao +  "|n: " + area.normal


        if (!area.st.isNullOrEmpty() && !area.numeroCorte.isNullOrEmpty()){
            var progresso: Float
            if (area.numeroCorte.toInt() < 1) {
                progresso = (area.st.toFloat() / ST_PRIRO_CORTE) * 100
            }else {
                progresso =  (area.st.toFloat() / ST_OUTROS_CORTES ) * 100
            }

            var pro = progresso.toInt()
            if (pro != null){
                tvProgresso?.text = pro.toString() + "%"
            }

        }else{
            tvProgresso?.text = "Progresso: ERRO"
        }


    }



}

